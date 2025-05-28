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
	"time"
)

// ClientHandler holds connection and user reference
type ClientHandler struct {
	Conn             net.Conn
	User             *User
	HandlerID        int
	StopMatchPending chan struct{}
}

// SendPDU sends a PDU to the server
func SendPDU(conn net.Conn, pdu PDU) error {
	data, err := json.Marshal(pdu)
	log.Println("Send data:", data)
	if err != nil {
		return fmt.Errorf("marshal error: %v", err)
	}

	// Send length prefix
	lenBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBuf, uint32(len(data)))
	log.Println("Send length:", lenBuf)
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
	if _, err := io.ReadFull(conn, lenBuf); err != nil {
		return PDU{}, fmt.Errorf("read length error: %w", err)
	}
	length := binary.BigEndian.Uint32(lenBuf)
	log.Printf("Received message of length %d", length)
	// Read PDU data
	data := make([]byte, length)
	if _, err := io.ReadFull(conn, data); err != nil {
		return PDU{}, fmt.Errorf("read data error: %w", err)
	}
	log.Printf("Received message %d", data)
	var pdu PDU
	if err := json.Unmarshal(data, &pdu); err != nil {
		return PDU{}, fmt.Errorf("unmarshal error: %w", err)
	}
	log.Printf("Received message data: %s", string(pdu.Data))
	log.Printf("Received message type: %s", string(pdu.Type))
	return pdu, nil
}

// HandleConnection manages a single client connection
func HandleConnection(conn net.Conn, users map[string]User, matchQueue chan *ClientHandler, id int) error {
	defer conn.Close()

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

	// Channel to signal when to stop sending match_pending
	stopPending := make(chan struct{})

	// Start goroutine to send match_pending every 5 seconds
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				SendPDU(conn, PDU{Type: "match_pending", Data: []byte(`{}`)})
			case <-stopPending:
				return
			}
		}
	}()

	// Enqueue for matchmaking
	handler := &ClientHandler{Conn: conn, User: &stored, HandlerID: id, StopMatchPending: stopPending}
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

	// Stop match_pending loops
	close(c1.StopMatchPending)
	close(c2.StopMatchPending)

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
