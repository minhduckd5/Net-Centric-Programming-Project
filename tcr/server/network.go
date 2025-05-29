// network.go
package server

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"tcr/specs"
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
	// log.Println("Send length:", lenBuf)
	if _, err := conn.Write(lenBuf); err != nil {
		log.Printf("write length error: %v", err)
		return fmt.Errorf("write length error: %w", err)
	}

	// Send PDU data
	if _, err := conn.Write(data); err != nil {
		return fmt.Errorf("write data error: %w", err)
	}
	// log.Println("Send data:", data)
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

func saveUsers(path string, users map[string]User) error {
	data, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// HandleConnection manages a single client connection
const userFilePath = "./players.json"

func HandleConnection(conn net.Conn, users map[string]User, matchQueue chan *ClientHandler, id int) {
	for {
		pdu, err := ReceivePDU(conn)
		if err != nil {
			log.Println("receive error:", err)
			return
		}

		var creds struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := json.Unmarshal(pdu.Data, &creds); err != nil {
			log.Println("unmarshal credentials:", err)
			return
		}

		switch pdu.Type {

		case "register":
			if _, exists := users[creds.Username]; exists {
				SendPDU(conn, PDU{
					Type: "register_resp",
					Data: []byte(`{"status":"ERR:UserExists"}`),
				})
				continue // ❗ Allow retry
			}

			newUser := User{
				Username:     creds.Username,
				PasswordHash: creds.Password,
				Level:        1,
				Exp:          0,
				NextLevel:    2,
				Multiplier:   2,
			}

			users[creds.Username] = newUser

			if err := saveUsers(userFilePath, users); err != nil {
				log.Println("error saving users:", err)
				SendPDU(conn, PDU{
					Type: "register_resp",
					Data: []byte(`{"status":"ERR:SaveFailed"}`),
				})
				continue // ❗ Allow retry
			}

			SendPDU(conn, PDU{
				Type: "register_resp",
				Data: []byte(`{"status":"OK"}`),
			})
			log.Printf("User registered: %s\n", creds.Username)

			// ✅ After registration, let them login in next loop
			continue

		case "login":
			stored, ok := users[creds.Username]
			if !ok || stored.PasswordHash != creds.Password {
				SendPDU(conn, PDU{
					Type: "login_resp",
					Data: []byte(`{"status":"ERR"}`),
				})
				continue // ❗ Allow retry
			}

			SendPDU(conn, PDU{
				Type: "login_resp",
				Data: []byte(`{"status":"OK"}`),
			})
			log.Printf("User logged in: %s\n", creds.Username)

			// ✅ Success: enqueue and exit loop
			handler := &ClientHandler{Conn: conn, User: &stored, HandlerID: id}
			matchQueue <- handler
			return

		default:
			SendPDU(conn, PDU{
				Type: "error",
				Data: []byte(`{"msg":"invalid command"}`),
			})
			continue
		}
	}
}

// StartServer begins listening and handles matchmaking
func StartServer(addr string, users map[string]User,
	troopSpecs map[string]specs.TroopSpec, towerSpecs map[string]specs.TowerSpec,
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
	troopSpecs map[string]specs.TroopSpec, towerSpecs map[string]specs.TowerSpec) {
	log.Println("accept:", c1, c2)
	// Send game_start PDU
	startData := fmt.Sprintf(`{"players":[%d,%d]}`, c1.HandlerID, c2.HandlerID)
	// fmt.Println("startData: ", startData)
	if err := SendPDU(c1.Conn, PDU{
		Type: "game_start",
		Data: []byte(startData)}); err != nil {
	}

	if err := SendPDU(c2.Conn, PDU{
		Type: "game_start",
		Data: []byte(startData)}); err != nil {
	}
	// fmt.Println("Send success")

	// Initialize session
	players := [2]*Player{
		{
			Conn:     c1.Conn,
			Username: c1.User.Username,
			Mana:     5,
			Towers: []*specs.TowerSpec{
				cloneTowerSpec(towerSpecs["guard_tower"]),
				cloneTowerSpec(towerSpecs["guard_tower"]),
				cloneTowerSpec(towerSpecs["king_tower"]),
			},
			Level: Level{
				Level:      c1.User.Level,
				Exp:        c1.User.Exp,
				NextLevel:  c1.User.NextLevel,
				Multiplier: c1.User.Multiplier,
			},
		},
		{
			Conn:     c2.Conn,
			Username: c2.User.Username,
			Mana:     5,
			Towers: []*specs.TowerSpec{
				cloneTowerSpec(towerSpecs["guard_tower"]),
				cloneTowerSpec(towerSpecs["guard_tower"]),
				cloneTowerSpec(towerSpecs["king_tower"]),
			},
			Level: Level{
				Level:      c2.User.Level,
				Exp:        c2.User.Exp,
				NextLevel:  c2.User.NextLevel,
				Multiplier: c2.User.Multiplier,
			},
		},
	}
	gs := NewGameSession(players, troopSpecs, towerSpecs)
	gs.StartGame()

}

func cloneTowerSpec(spec specs.TowerSpec) *specs.TowerSpec {
	clone := spec // copy struct value
	return &clone
}
