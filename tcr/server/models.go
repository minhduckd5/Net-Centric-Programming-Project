// server/models.go
package main

import (
	"encoding/json"
	"os"
)

// User represents a player account
type User struct {
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
	Level        int    `json:"level"`
	EXP          int    `json:"exp"`
}

// // TroopSpec loaded from specs.json
// type TroopSpec struct {
// 	HP   int `json:"hp"`
// 	ATK  int `json:"atk"`
// 	DEF  int `json:"def"`
// 	Mana int `json:"mana"`
// 	EXP  int `json:"exp"`
// }

// TowerSpec loaded from specs.json
type TowerSpec struct {
	HP   int     `json:"hp"`
	ATK  int     `json:"atk"`
	DEF  int     `json:"def"`
	Crit float64 `json:"crit"`
	EXP  int     `json:"exp"`
}

// PDU is the envelope for all messages (newline‑delimited JSON)
type PDU struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

// LoadUsers reads the users JSON file
func LoadUsers(path string) (map[string]User, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var users map[string]User
	if err := json.NewDecoder(file).Decode(&users); err != nil {
		return nil, err
	}
	return users, nil
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
