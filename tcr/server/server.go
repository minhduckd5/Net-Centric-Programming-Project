// server.go
package server

import (
	// "tcr/common"

	"tcr/config"
	"tcr/specs"
)

// GameManager handles all active game sessions
type GameManager struct {
	sessions   map[string]*GameSession
	matchQueue chan *ClientHandler
	// mutex      sync.RWMutex
	specs  *specs.Specs
	config *config.Config
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
// func convertSpecs(s *specs.Specs) (map[string]specs.TroopSpec, map[string]specs.TowerSpec) {
// 	troopSpecs := make(map[string]specs.TroopSpec)
// 	towerSpecs := make(map[string]specs.TowerSpec)

// 	for k, v := range s.Troops {
// 		troopSpecs[k] = specs.TroopSpec{
// 			Name:        v.Name,
// 			Type:        v.Type,
// 			Health:      v.Health,
// 			Damage:      v.Damage,
// 			Defence:     v.Defence,
// 			Range:       v.Range,
// 			Speed:       v.Speed,
// 			AttackSpeed: v.AttackSpeed,
// 			Cost:        v.Cost,
// 			Target:      v.Target,
// 		}
// 	}

// 	for k, v := range s.Towers {
// 		towerSpecs[k] = specs.TowerSpec{
// 			Name:        v.Name,
// 			Type:        v.Type,
// 			Health:      v.Health,
// 			Damage:      v.Damage,
// 			Defence:     v.Defence,
// 			Range:       v.Range,
// 			AttackSpeed: v.AttackSpeed,
// 			Target:      v.Target,
// 		}
// 	}

// 	return troopSpecs, towerSpecs
// }

// StartMatchmaking starts the matchmaking process
// func (gm *GameManager) StartMatchmaking() {
// 	go func() {
// 		for {
// 			// Wait for two players
// 			player1 := <-gm.matchQueue
// 			player2 := <-gm.matchQueue

// 			// Create new game session
// 			sessionID := fmt.Sprintf("game_%d", time.Now().UnixNano())
// 			players := [2]*Player{
// 				{Conn: player1.Conn, Username: player1.User.Username},
// 				{Conn: player2.Conn, Username: player2.User.Username},
// 			}

// 			// Convert specs to game types
// 			troopSpecs, towerSpecs := convertSpecs(gm.specs)

// 			// Initialize game session
// 			session := NewGameSession(players, troopSpecs, towerSpecs)

// 			// Store session
// 			gm.mutex.Lock()
// 			gm.sessions[sessionID] = session
// 			gm.mutex.Unlock()

// 			// Start game in goroutine
// 			go func() {
// 				session.StartGame()
// 				// Clean up when game ends
// 				gm.mutex.Lock()
// 				delete(gm.sessions, sessionID)
// 				gm.mutex.Unlock()
// 			}()
// 		}
// 	}()
// }
