// server/models.go
package main

import (
	"encoding/json"
	"net"
	"os"
	"time"
)

// User represents a player account persisted in players.json
type User struct {
	Username     string         `json:"username"`
	PasswordHash string         `json:"password_hash"`
	Level        int            `json:"level"`
	EXP          int            `json:"exp"`
	TroopLevels  map[string]int `json:"troop_levels"`
	TowerLevels  map[string]int `json:"tower_levels"`
}

// TroopSpec loaded from specs.json
type TroopSpec struct {
	Name string
	ATK  int
	DEF  int
	Mana int
}

// TowerSpec loaded from specs.json
type TowerSpec struct {
	Name string
	HP   int
	DEF  int
}

// PDU is the envelope for all messages (newline‑delimited JSON)
type PDU struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

// Player represents a game player
type Player struct {
	Username string
	Password string
	Mana     int
	Towers   []*Tower
	Conn     net.Conn
}

// Tower represents a game tower
type Tower struct {
	Name string
	HP   int
	DEF  int
}

// GameSession holds state for a single 1v1 match
type GameSession struct {
	Mode               GameMode
	Players            [2]*Player
	TroopSpecs         map[string]TroopSpec
	TowerSpecs         map[string]TowerSpec
	Commands           chan DeployCmd
	Done               chan struct{}
	TickInterval       time.Duration
	justDestroyedTower bool // tracks if a tower was just destroyed
}

// DeployCmd is issued by a client or AI to deploy a troop
type DeployCmd struct {
	PlayerIndex int
	TroopName   string
}

// GameMode indicates simple (turn-based) or enhanced (continuous)
type GameMode int

const (
	SimpleMode GameMode = iota
	EnhancedMode
)

// LoadUsers reads players.json into a map
func LoadUsers(path string) (map[string]User, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var m map[string]User
	if err := json.NewDecoder(file).Decode(&m); err != nil {
		return nil, err
	}
	return m, nil
}

// SaveUsers writes the user map back to JSON file
func SaveUsers(path string, m map[string]User) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	return enc.Encode(m)
}
