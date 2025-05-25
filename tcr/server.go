// server.go
package main

import (
	// "tcr/common"

	"fmt"
	"sync"
	"time"

	"tcr/config"
	"tcr/specs"
)

// GameManager handles all active game sessions
type GameManager struct {
	sessions   map[string]*GameSession
	matchQueue chan *ClientHandler
	mutex      sync.RWMutex
	specs      *specs.Specs
	config     *config.Config
}

// NewGameManager creates a new game manager
func NewGameManager(specs *specs.Specs, config *config.Config) *GameManager {
	return &GameManager{
		sessions:   make(map[string]*GameSession),
		matchQueue: make(chan *ClientHandler, config.Game.MaxPlayers),
		specs:      specs,
		config:     config,
	}
}

// convertSpecs converts specs package types to game types
func convertSpecs(specs *specs.Specs) (map[string]TroopSpec, map[string]TowerSpec) {
	troopSpecs := make(map[string]TroopSpec)
	towerSpecs := make(map[string]TowerSpec)

	for k, v := range specs.Troops {
		troopSpecs[k] = TroopSpec{
			HP:   v.Health,
			ATK:  v.Damage,
			DEF:  0, // Default defense
			Mana: v.Cost,
			EXP:  0, // Default experience
		}
	}

	for k, v := range specs.Towers {
		towerSpecs[k] = TowerSpec{
			HP:   v.Health,
			ATK:  v.Damage,
			DEF:  0,   // Default defense
			Crit: 0.1, // Default crit chance
			EXP:  0,   // Default experience
		}
	}

	return troopSpecs, towerSpecs
}

// StartMatchmaking starts the matchmaking process
func (gm *GameManager) StartMatchmaking() {
	go func() {
		for {
			// Wait for two players
			player1 := <-gm.matchQueue
			player2 := <-gm.matchQueue

			// Create new game session
			sessionID := fmt.Sprintf("game_%d", time.Now().UnixNano())
			players := [2]*Player{
				{Conn: player1.Conn, Username: player1.User.Username},
				{Conn: player2.Conn, Username: player2.User.Username},
			}

			// Convert specs to game types
			troopSpecs, towerSpecs := convertSpecs(gm.specs)

			// Initialize game session
			session := NewGameSession(EnhancedMode, players, troopSpecs, towerSpecs)

			// Store session
			gm.mutex.Lock()
			gm.sessions[sessionID] = session
			gm.mutex.Unlock()

			// Start game in goroutine
			go func() {
				session.StartGame()
				// Clean up when game ends
				gm.mutex.Lock()
				delete(gm.sessions, sessionID)
				gm.mutex.Unlock()
			}()
		}
	}()
}
