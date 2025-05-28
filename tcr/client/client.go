// client/main.go
package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
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
	serverAddr string
	conn       net.Conn
	reader     *bufio.Reader
	username   string
	password   string
	inGame     bool
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
	fmt.Printf("Your Mana: %d\n", state.YourMana)
	fmt.Printf("Opponent Mana: %d\n", state.OpponentMana)

	fmt.Println("\nYour Towers:")
	for _, tower := range state.YourTowers {
		fmt.Printf("- %s: HP %d\n", tower.Name, tower.HP)
	}

	fmt.Println("\nOpponent Towers:")
	for _, tower := range state.OpponentTowers {
		fmt.Printf("- %s: HP %d\n", tower.Name, tower.HP)
	}

	fmt.Println("\nAvailable Troops:")
	fmt.Println("1. Pawn (3 mana)")
	fmt.Println("2. Bishop (4 mana)")
	fmt.Println("3. Rook (5 mana)")
	fmt.Println("4. Knight (5 mana)")
	fmt.Println("5. Prince (6 mana)")
	fmt.Println("6. Queen (5 mana)")
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

	// Main game loop
	for {
		pdu, err := server.ReceivePDU(c.conn)
		if err != nil {
			if errors.Is(err, io.EOF) {
				// Now this will correctly catch wrapped EOFs
				fmt.Println("Waiting for server message...")
				time.Sleep(1 * time.Second)
				continue
			}
			if c.inGame {
				fmt.Printf("Connection lost. Attempting to reconnect...\n")
				if err := c.connect(); err != nil {
					return fmt.Errorf("reconnection failed: %v", err)
				}
				continue
			}
			return fmt.Errorf("receive error: %v", err)
		}

		switch pdu.Type {
		case "game_start":
			c.handleGameStart(pdu)
		case "state_update":
			c.handleStateUpdate(pdu)
		case "game_end":
			c.handleGameEnd(pdu)
			return nil
		}

		if c.inGame {
			fmt.Print("\nEnter troop number (1-6) or 'quit': ")
			input := strings.TrimSpace(readLine(c.reader))

			if input == "quit" {
				return nil
			}

			// Map input to troop names
			troopMap := map[string]string{
				"1": "Pawn",
				"2": "Bishop",
				"3": "Rook",
				"4": "Knight",
				"5": "Prince",
				"6": "Queen",
			}

			if troop, ok := troopMap[input]; ok {
				deployCmd := fmt.Sprintf(`{"troop":"%s"}`, troop)
				if err := server.SendPDU(c.conn, server.PDU{Type: "deploy", Data: json.RawMessage(deployCmd)}); err != nil {
					fmt.Printf("Error sending deploy command: %v\n", err)
				}
			} else {
				fmt.Println("Invalid troop number!")
			}
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
