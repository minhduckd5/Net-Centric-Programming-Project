// network.go
package server

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
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
		return fmt.Errorf("marshal error: %w", err)
	}

	// Send length prefix
	lenBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBuf, uint32(len(data)))
	log.Println("Send length:", lenBuf)
	if _, err := conn.Write(lenBuf); err != nil {
		log.Printf("write length error: %v", err)
		return fmt.Errorf("write length error: %w", err)
	}

	// Send PDU data
	if _, err := conn.Write(data); err != nil {
		return fmt.Errorf("write data error: %w", err)
	}
	log.Println("Send data:", data)
	return nil
}

// ReceivePDU receives a PDU from the server
func ReceivePDU(conn net.Conn) (PDU, error) {
	// Read length prefix
	lenBuf := make([]byte, 4)
	if _, err := io.ReadFull(conn, lenBuf); err != nil {
		return PDU{}, fmt.Errorf("read length error: %w", err)
	}
	length := binary.BigEndian.Uint32(lenBuf)
	// log.Printf("Received message of length %d", length)
	// Read PDU data
	data := make([]byte, length)
	if _, err := io.ReadFull(conn, data); err != nil {
		return PDU{}, fmt.Errorf("read data error: %w", err)
	}
	// log.Printf("Received message %d", data)
	var pdu PDU
	if err := json.Unmarshal(data, &pdu); err != nil {
		return PDU{}, fmt.Errorf("unmarshal error: %w", err)
	}
	// log.Printf("Received message data: %s", string(pdu.Data))
	// log.Printf("Received message type: %s", string(pdu.Type))
	return pdu, nil
}

// HandleConnection manages a single client connection
func HandleConnection(conn net.Conn, users map[string]User, matchQueue chan *ClientHandler, id int) error {

	pdu, err := ReceivePDU(conn)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}

	var creds struct{ Username, Password string }
	if err := json.Unmarshal(pdu.Data, &creds); err != nil {
		log.Println("unmarshal credentials:", err)
		return fmt.Errorf("error: %v", err)
	}
	log.Printf("Received login credentials: %s / %s", creds.Username, creds.Password)

	writer := bufio.NewWriter(conn)
	enc := json.NewEncoder(writer)

	stored, ok := users[creds.Username]
	if !ok || stored.PasswordHash != creds.Password {
		enc.Encode(PDU{Type: "login_resp", Data: []byte(`{"status":"ERR"}`)})
		writer.Flush()
		return fmt.Errorf("invalid login for user: %s", creds.Username)
	}

	// Send login success
	resp := PDU{Type: "login_resp", Data: []byte(`{"status":"OK"}`)}
	if err := SendPDU(conn, resp); err != nil {
		log.Println("send login_resp:", err)
		return err
	}

	// Enqueue for matchmaking
	handler := &ClientHandler{Conn: conn, User: &stored, HandlerID: id}
	matchQueue <- handler

	return nil
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
		for {
			log.Println("Waiting for first client...")
			c1 := <-matchQueue
			log.Printf("Got first client: %v", c1.User.Username)

			log.Println("Waiting for second client...")
			c2 := <-matchQueue
			log.Printf("Got second client: %v", c2.User.Username)

			log.Println("Starting game session...")
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
	troopSpecs map[string]TroopSpec, towerSpecs map[string]TowerSpec) error {
	log.Println("accept:", c1, c2)
	// Send game_start PDU
	startData := fmt.Sprintf(`{"mode":"simple","players":[%d,%d]}`, c1.HandlerID, c2.HandlerID)
	fmt.Println("startData: ", startData)
	if err := SendPDU(c1.Conn, PDU{
		Type: "game_start",
		Data: []byte(startData)}); err != nil {
		return fmt.Errorf("login send error: %w", err)
	}

	if err := SendPDU(c2.Conn, PDU{
		Type: "game_start",
		Data: []byte(startData)}); err != nil {
		return fmt.Errorf("login send error: %w", err)
	}
	fmt.Printf("Send success")

	// Initialize session
	players := [2]*Player{
		{Conn: c1.Conn, Username: c1.User.Username},
		{Conn: c2.Conn, Username: c2.User.Username},
	}
	gs := NewGameSession(SimpleMode, players, troopSpecs, towerSpecs)
	gs.StartGame()
	return nil
}
