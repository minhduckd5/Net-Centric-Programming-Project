// network.go
package server

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net"
)

// ClientHandler holds connection and user reference
type ClientHandler struct {
	Conn      net.Conn
	User      *User
	HandlerID int
}

// SendPDU sends a PDU to the server
func SendPDU(conn net.Conn, pdu PDU) error {
	data, err := json.Marshal(pdu)
	if err != nil {
		return fmt.Errorf("marshal error: %v", err)
	}

	// Send length prefix
	lenBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBuf, uint32(len(data)))
	if _, err := conn.Write(lenBuf); err != nil {
		return fmt.Errorf("write length error: %v", err)
	}

	// Send PDU data
	if _, err := conn.Write(data); err != nil {
		return fmt.Errorf("write data error: %v", err)
	}

	return nil
}

// ReceivePDU receives a PDU from the server
func ReceivePDU(conn net.Conn) (PDU, error) {
	// Read length prefix
	lenBuf := make([]byte, 4)
	if _, err := conn.Read(lenBuf); err != nil {
		return PDU{}, fmt.Errorf("read length error: %v", err)
	}
	length := binary.BigEndian.Uint32(lenBuf)

	// Read PDU data
	data := make([]byte, length)
	if _, err := conn.Read(data); err != nil {
		return PDU{}, fmt.Errorf("read data error: %v", err)
	}

	var pdu PDU
	if err := json.Unmarshal(data, &pdu); err != nil {
		return PDU{}, fmt.Errorf("unmarshal error: %v", err)
	}

	return pdu, nil
}

// HandleConnection manages a single client connection
func HandleConnection(conn net.Conn, users map[string]User, matchQueue chan *ClientHandler, id int) {
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

	// Credential check
	var creds struct{ Username, Password string }
	json.Unmarshal(pdu.Data, &creds)
	stored, ok := users[creds.Username]
	if !ok || stored.PasswordHash != creds.Password {
		enc.Encode(PDU{Type: "login_resp", Data: []byte(`{"status":"ERR"}`)})
		writer.Flush()
		return
	}

	// Success
	enc.Encode(PDU{Type: "login_resp", Data: []byte(`{"status":"OK"}`)})
	writer.Flush()

	// Enqueue for matchmaking
	matchQueue <- &ClientHandler{Conn: conn, User: &stored, HandlerID: id}
}

// StartServer begins listening and handles matchmaking
func StartServer(addr string, users map[string]User,
	troopSpecs map[string]TroopSpec, towerSpecs map[string]TowerSpec,
	matchQueue chan *ClientHandler) error {

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	log.Println("Server listening on", addr)

	go func() {
		// Matchmaking goroutine: pair two clients
		for {
			c1 := <-matchQueue
			c2 := <-matchQueue
			go StartGameSession(c1, c2, troopSpecs, towerSpecs)
		}
	}()

	handlerID := 0
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("accept:", err)
			continue
		}
		handlerID++
		go HandleConnection(conn, users, matchQueue, handlerID)
	}
}

// StartGameSession initializes GameSession and triggers startGame
func StartGameSession(c1, c2 *ClientHandler,
	troopSpecs map[string]TroopSpec, towerSpecs map[string]TowerSpec) {

	// Send game_start PDU
	startData := fmt.Sprintf(`{"mode":"simple","players":[%d,%d]}`, c1.HandlerID, c2.HandlerID)
	SendPDU(c1.Conn, PDU{Type: "game_start", Data: []byte(startData)})
	SendPDU(c2.Conn, PDU{Type: "game_start", Data: []byte(startData)})

	// Initialize session
	players := [2]*Player{
		{Conn: c1.Conn, Username: c1.User.Username},
		{Conn: c2.Conn, Username: c2.User.Username},
	}
	gs := NewGameSession(SimpleMode, players, troopSpecs, towerSpecs)
	gs.StartGame()
}
