// models.go
package server

import (
	"encoding/json"
	"net"
	"os"
	"tcr/specs"
)

// User represents a player account
type User struct {
	Username     string  `json:"username"`
	PasswordHash string  `json:"password_hash"`
	Level        int     `json:"level"`
	Exp          int     `json:"exp"`
	NextLevel    int     `json:"next_level"`
	Multiplier   float64 `json:"multiplier"`
}

// Player represents a player in a game session
type Player struct {
	Conn        net.Conn
	Username    string
	Mana        int
	Towers      []*specs.TowerSpec
	Level       Level
	ActiveTroops []*TroopInstance // Or a similar struct you define
}

// PDU represents a Protocol Data Unit for client-server communication
type PDU struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

// Methods for Player
func (p *Player) NextAliveTower() *specs.TowerSpec {
	for _, t := range p.Towers {
		if t.Health > 0 {
			return t
		}
	}
	return nil
}

func (p *Player) DestroyTower(t *specs.TowerSpec) {
	t.Health = 0
}

func (p *Player) KingTowerDestroyed() bool {
	for _, t := range p.Towers {
		if t.Name == "King Tower" && t.Health <= 0 {
			return true
		}
	}
	return false
}

func (p *Player) HealWeakestTower(amount int) {
	var weakest *specs.TowerSpec
	minHP := 999999
	for _, t := range p.Towers {
		if t.Health > 0 && t.Health < minHP {
			weakest = t
			minHP = t.Health
		}
	}
	if weakest != nil {
		weakest.Health += amount
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
