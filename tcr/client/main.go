// client/main.go
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
)

func main() {
	addr := flag.String("addr", "localhost:9000", "server address")
	flag.Parse()

	conn, err := net.Dial("tcp", *addr)
	if err != nil {
		fmt.Println("connect error:", err)
		return
	}
	defer conn.Close()

	// Prompt for credentials
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Username: ")
	user, _ := reader.ReadString('\n')
	fmt.Print("Password: ")
	pass, _ := reader.ReadString('\n')

	// Send login PDU
	cred := fmt.Sprintf(`{"username":"%s","password":"%s"}`,
		user[:len(user)-1], pass[:len(pass)-1])
	SendPDU(conn, PDU{Type: "login", Data: json.RawMessage(cred)})

	// Await login_resp
	pdu, err := ReceivePDU(conn)
	if err != nil {
		fmt.Println("recv error:", err)
		return
	}
	fmt.Println("Server:", string(pdu.Data))

	// TODO: handle game_start, then loop for deploy/state_update/game_end
	for {
		pdu, err := ReceivePDU(conn)
		if err != nil {
			fmt.Println("recv error:", err)
			return
		}

		switch pdu.Type {
		case "game_start":
			fmt.Println("Game started!")
			fmt.Println(string(pdu.Data))
		case "state_update":
			var state GameState
			if err := json.Unmarshal(pdu.Data, &state); err != nil {
				fmt.Println("state parse error:", err)
				continue
			}
			displayGameState(state)
		case "game_end":
			fmt.Println("Game Over!")
			fmt.Println(string(pdu.Data))
			return
		}

		// Handle user input for troop deployment
		fmt.Print("\nEnter troop to deploy (or 'quit' to exit): ")
		input, _ := reader.ReadString('\n')
		input = input[:len(input)-1]

		if input == "quit" {
			return
		}

		// Send deploy command
		deployCmd := fmt.Sprintf(`{"troop":"%s"}`, input)
		SendPDU(conn, PDU{Type: "deploy", Data: json.RawMessage(deployCmd)})
	}
}

func displayGameState(state GameState) {
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
	fmt.Println("================")
}
