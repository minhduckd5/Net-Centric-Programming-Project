// models.go
package main

import (
	"encoding/json"
	"net"
	"os"
)

// User represents a player account
type User struct {
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
	Experience   int    `json:"experience"`
}

// Player represents a player in a game session
type Player struct {
	Conn     net.Conn
	Username string
	Mana     int
	Towers   []*Tower
	Level    Level
}

// Tower represents a game tower
type Tower struct {
	Name string `json:"name"`
	HP   int    `json:"hp"`
	DEF  int    `json:"def"`
	ATK  int    `json:"atk"`
}

// TroopSpec defines the stats for a troop type
type TroopSpec struct {
	HP   int `json:"hp"`
	ATK  int `json:"atk"`
	DEF  int `json:"def"`
	Mana int `json:"mana"`
	EXP  int `json:"exp"`
}

// TowerSpec defines the stats for a tower type
type TowerSpec struct {
	HP   int     `json:"hp"`
	ATK  int     `json:"atk"`
	DEF  int     `json:"def"`
	Crit float64 `json:"crit"`
	EXP  int     `json:"exp"`
}

// PDU represents a Protocol Data Unit for client-server communication
type PDU struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

// Methods for Player
func (p *Player) NextAliveTower() *Tower {
	for _, t := range p.Towers {
		if t.HP > 0 {
			return t
		}
	}
	return nil
}

func (p *Player) DestroyTower(t *Tower) {
	t.HP = 0
}

func (p *Player) KingTowerDestroyed() bool {
	for _, t := range p.Towers {
		if t.Name == "King Tower" && t.HP <= 0 {
			return true
		}
	}
	return false
}

func (p *Player) HealWeakestTower(amount int) {
	var weakest *Tower
	minHP := 999999
	for _, t := range p.Towers {
		if t.HP > 0 && t.HP < minHP {
			weakest = t
			minHP = t.HP
		}
	}
	if weakest != nil {
		weakest.HP += amount
	}
}

// LoadUsers reads user data from a JSON file
func LoadUsers(filename string) (map[string]User, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var users map[string]User
	if err := json.Unmarshal(data, &users); err != nil {
		return nil, err
	}
	return users, nil
}

// LoadSpecs reads troop and tower specifications from a JSON file
func LoadSpecs(filename string, troops map[string]TroopSpec, towers map[string]TowerSpec) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	var specs struct {
		Troops map[string]TroopSpec `json:"troops"`
		Towers map[string]TowerSpec `json:"towers"`
	}
	if err := json.Unmarshal(data, &specs); err != nil {
		return err
	}
	for k, v := range specs.Troops {
		troops[k] = v
	}
	for k, v := range specs.Towers {
		towers[k] = v
	}
	return nil
}
