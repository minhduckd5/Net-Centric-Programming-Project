// server/network.go
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
)

// type PDU struct {
// 	Type string          `json:"type"`
// 	Data json.RawMessage `json:"data"`
// }

// type User struct {
// 	// Define User fields as needed; for example:
// 	Username string `json:"username"`
// 	Password string `json:"password"`
// }

// handleConnection manages a single client
func handleConnection(conn net.Conn, users map[string]User) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	enc := json.NewEncoder(writer)

	// 1) AUTHENTICATION PHASE
	line, err := reader.ReadString('\n')
	if err != nil {
		log.Println("read login:", err)
		return
	}
	var pdu PDU
	if err := json.Unmarshal([]byte(line), &pdu); err != nil {
		log.Println("unmarshal pdu:", err)
		return
	}
	if pdu.Type != "login" {
		log.Println("expected login PDU")
		return
	}
	// TODO: validate credentials in pdu.Data, respond with login_resp PDU

	// 2) MATCHMAKING (stub)
	// Once two players connect, call startGame(session)

	// 3) GAME LOOP (in game.go)

	// Example: send a login_resp
	resp := PDU{Type: "login_resp", Data: json.RawMessage(`{"status":"OK"}`)}
	if err := enc.Encode(resp); err != nil {
		log.Println("send login_resp:", err)
	}
	writer.Flush()
}

// StartServer initializes the game server
func StartServer(addr, usersPath string) error {
	// Load user data
	users, err := LoadUsers(usersPath)
	if err != nil {
		return fmt.Errorf("load users: %v", err)
	}

	// Start TCP listener
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("listen: %v", err)
	}
	defer listener.Close()

	log.Printf("Server listening on %s", addr)

	// Accept connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("accept error: %v", err)
			continue
		}
		go handleConnection(conn, users)
	}
}
