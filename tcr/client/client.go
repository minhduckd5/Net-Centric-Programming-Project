// client/main.go
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"tcr/server"
	"time"
)

const (
	maxReconnectAttempts = 3
	reconnectDelay       = 5 * time.Second
)

type GameClient struct {
	serverAddr      string
	conn            net.Conn
	reader          *bufio.Reader
	username        string
	password        string
	inGame          bool
	availableTroops []string
}

func NewGameClient(serverAddr string) *GameClient {
	return &GameClient{
		serverAddr: serverAddr,
		reader:     bufio.NewReader(os.Stdin),
	}
}

func (c *GameClient) connect() error {
	var err error
	for i := 0; i < maxReconnectAttempts; i++ {
		c.conn, err = net.Dial("tcp", c.serverAddr)
		if err == nil {
			return nil
		}
		fmt.Printf("Connection attempt %d failed: %v\n", i+1, err)
		if i < maxReconnectAttempts-1 {
			fmt.Printf("Retrying in %v...\n", reconnectDelay)
			time.Sleep(reconnectDelay)
		}
	}
	return fmt.Errorf("failed to connect after %d attempts", maxReconnectAttempts)
}

func (c *GameClient) login() error {
	fmt.Print("Username: ")
	c.username = strings.TrimSpace(readLine(c.reader))
	fmt.Print("Password: ")
	c.password = strings.TrimSpace(readLine(c.reader))

	cred := fmt.Sprintf(`{"username":"%s","password":"%s"}`, c.username, c.password)
	if err := server.SendPDU(c.conn, server.PDU{
		Type: "login",
		Data: json.RawMessage(cred)}); err != nil {
		return fmt.Errorf("login send error: %v", err)
	}

	pdu, err := server.ReceivePDU(c.conn)
	if err != nil {
		return fmt.Errorf("login response error: %v", err)
	}

	var resp struct{ Status string }
	if err := json.Unmarshal(pdu.Data, &resp); err != nil {
		return fmt.Errorf("login parse error: %v", err)
	}

	if resp.Status != "OK" {
		return fmt.Errorf("login failed: %s", resp.Status)
	}

	fmt.Println("Login successful!")
	return nil
}

func (c *GameClient) handleGameStart(pdu server.PDU) {
	var startData struct {
		Mode    string `json:"mode"`
		Players []int  `json:"players"`
	}
	if err := json.Unmarshal(pdu.Data, &startData); err != nil {
		fmt.Printf("Error parsing game start: %v\n", err)
		return
	}

	// Randomly select 3 unique troops at game start
	allTroops := []string{"pawn", "bishop", "rook", "knight", "prince", "queen", "archer", "giant", "minion"}
	rand.Shuffle(len(allTroops), func(i, j int) { allTroops[i], allTroops[j] = allTroops[j], allTroops[i] })
	c.availableTroops = allTroops[:3]
	c.inGame = true
	fmt.Printf("\n=== Game Started ===\n")
	fmt.Printf("Mode: %s\n", startData.Mode)
	fmt.Printf("Players: %v\n", startData.Players)
	fmt.Println("==================")
}

func (c *GameClient) handleStateUpdate(pdu server.PDU) {
	var state server.GameState
	if err := json.Unmarshal(pdu.Data, &state); err != nil {
		fmt.Printf("Error parsing state update: %v\n", err)
		return
	}

	// Clear screen
	fmt.Print("\033[H\033[2J")

	// Display game state
	fmt.Println("\n=== Game State ===")
	fmt.Printf("Player 1 Mana: %d\n", state.YourMana)
	fmt.Printf("Player 2 Mana: %d\n", state.OpponentMana)

	fmt.Println("\nPlayer 1 Towers:")
	for _, tower := range state.Player1Towers {
		fmt.Printf("- %s: HP %d\n", tower.Name, tower.Health)
	}

	fmt.Println("\nPlayer 2 Towers:")
	for _, tower := range state.Player2Towers {
		fmt.Printf("- %s: HP %d\n", tower.Name, tower.Health)
	}

	troopCosts := map[string]int{
		"pawn": 3, "bishop": 4, "rook": 5,
		"knight": 5, "prince": 6, "queen": 5,
		"archer": 3, "giant": 7, "minion": 3,
	}

	fmt.Println("\nAvailable Troops:")
	for i, troop := range c.availableTroops {
		capitalized := strings.Title(troop)
		fmt.Printf("%d. %s (%d mana)\n", i+1, capitalized, troopCosts[troop])
	}

	fmt.Println("\nEnter troop number or 'quit' to exit")
}

func (c *GameClient) handleGameEnd(pdu server.PDU) {
	var endData struct {
		Result string `json:"result"`
		Exp    int    `json:"exp"`
	}
	if err := json.Unmarshal(pdu.Data, &endData); err != nil {
		fmt.Printf("Error parsing game end: %v\n", err)
		return
	}

	fmt.Printf("\n=== Game Over ===\n")
	fmt.Printf("Result: %s\n", endData.Result)
	fmt.Printf("EXP Gained: %d\n", endData.Exp)
	fmt.Println("================")

	c.inGame = false
}

func (c *GameClient) run() error {
	if err := c.connect(); err != nil {
		return err
	}
	defer c.conn.Close()

	if err := c.login(); err != nil {
		return err
	}

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\nShutting down...")
		if c.conn != nil {
			c.conn.Close()
		}
		os.Exit(0)
	}()

	// Start goroutine to receive updates
	go func() {
		for {
			pdu, err := server.ReceivePDU(c.conn)
			if err != nil {
				fmt.Printf("Receive error: %v\n", err)
				continue
			}

			switch pdu.Type {
			case "game_start":
				c.handleGameStart(pdu)
			case "state_update":
				c.handleStateUpdate(pdu)
			case "game_end":
				c.handleGameEnd(pdu)
				os.Exit(0) // Gracefully exit game
			}
		}
	}()

	// Input loop
	for {
		if c.inGame {
			fmt.Print("\nEnter troop number (1–3) or 'quit': ")
			input := strings.TrimSpace(readLine(c.reader))

			if input == "quit" {
				return nil
			}

			if idx, err := strconv.Atoi(input); err == nil && idx >= 1 && idx <= 3 {
				troop := c.availableTroops[idx-1]
				deployCmd := fmt.Sprintf(`{"Troop": "%s"}`, troop)
				if err := server.SendPDU(c.conn, server.PDU{
					Type: "deploy", Data: json.RawMessage(deployCmd)}); err != nil {
					fmt.Printf("Error sending deploy command: %v\n", err)
				}
			} else {
				fmt.Println("Invalid troop number!")
			}
		} else {
			time.Sleep(500 * time.Millisecond) // Avoid busy-waiting
		}
	}
}

func readLine(reader *bufio.Reader) string {
	line, _ := reader.ReadString('\n')
	return strings.TrimSpace(line)
}

func main() {
	serverAddr := flag.String("server", "localhost:9000", "Server address")
	flag.Parse()

	client := NewGameClient(*serverAddr)
	if err := client.run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
