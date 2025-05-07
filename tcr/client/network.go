// client/network.go
package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
)

// PDU represents a Protocol Data Unit for client-server communication
// type PDU struct {
// 	Type string          `json:"type"`
// 	Data json.RawMessage `json:"data"`
// }

// GameState represents the current state of the game
type GameState struct {
	YourMana       int     `json:"your_mana"`
	OpponentMana   int     `json:"opponent_mana"`
	YourTowers     []Tower `json:"your_towers"`
	OpponentTowers []Tower `json:"opponent_towers"`
}

// Tower represents a game tower
// type Tower struct {
// 	Name string `json:"name"`
// 	HP   int    `json:"hp"`
// }

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
