// game.go
// Core game logic for Text-Based Clash Royale (TCR)

package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// GameMode indicates simple (turn-based) or enhanced (continuous)
type GameMode int

const (
	SimpleMode GameMode = iota
	EnhancedMode
)

// GameSession holds state for a single 1v1 match
var mutex sync.Mutex // protects session state during concurrent access

// Level represents a player's level and associated stats
type Level struct {
	Level      int     `json:"level"`
	Exp        int     `json:"exp"`
	NextLevel  int     `json:"next_level"`
	Multiplier float64 `json:"multiplier"`
}

// Update Player struct
// type Player struct {
// 	Conn     net.Conn
// 	Username string
// 	Mana     int
// 	Towers   []*Tower
// 	Level    Level // Add level information
// }

type GameSession struct {
	Mode               GameMode
	Players            [2]*Player // two players
	TroopSpecs         map[string]TroopSpec
	TowerSpecs         map[string]TowerSpec
	Commands           chan DeployCmd // incoming deploy commands
	Done               chan struct{}  // signals end of game
	TickInterval       time.Duration  // for enhanced mode
	justDestroyedTower bool           // tracks if a tower was just destroyed
}

// DeployCmd is issued by a client or AI to deploy a troop
type DeployCmd struct {
	PlayerIndex int    // 0 or 1
	TroopName   string // e.g., "Pawn"
}

// GameState represents the current state of the game
type GameState struct {
	YourMana       int     `json:"your_mana"`
	OpponentMana   int     `json:"opponent_mana"`
	YourTowers     []Tower `json:"your_towers"`
	OpponentTowers []Tower `json:"opponent_towers"`
}

// startGame launches the appropriate game loop based on mode
func (gs *GameSession) StartGame() {
	switch gs.Mode {
	case SimpleMode:
		gs.simpleLoop()
	case EnhancedMode:
		gs.enhancedLoop()
	}
}

// simpleLoop runs turn-based gameplay
func (gs *GameSession) simpleLoop() {
	current := rand.Intn(2) // random start; seed earlier
	for {
		select {
		case cmd := <-gs.Commands:
			if cmd.PlayerIndex != current {
				// ignore commands out of turn
				continue
			}
			// process deploy
			gs.handleDeploy(cmd)
			// if no tower destroyed, switch turn
			if !gs.justDestroyedTower {
				current = 1 - current
			}
			// check win condition
			if gs.checkGameEnd() {
				close(gs.Done)
				return
			}
		case <-gs.Done:
			return
		}
	}
}

// enhancedLoop runs real-time gameplay with mana regen and timeout
func (gs *GameSession) enhancedLoop() {
	ticker := time.NewTicker(gs.TickInterval)
	timeout := time.After(3 * time.Minute) // 3-minute match timer
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			gs.tick() // regen mana, tower attacks, send state
		case cmd := <-gs.Commands:
			gs.handleDeploy(cmd) // immediate deploy handling
		case <-timeout:
			gs.evaluateWinner()
			close(gs.Done)
			return
		case <-gs.Done:
			return
		}
	}
}

// tick handles periodic updates: mana regen and optional tower attacks
func (gs *GameSession) tick() {
	mutex.Lock()
	defer mutex.Unlock()
	for i := range gs.Players {
		p := gs.Players[i]
		if p.Mana < 10 {
			p.Mana++ // mana regen
		}
	}
	// optionally process tower auto-attacks here
	gs.broadcastState()
}

// handleDeploy processes a DeployCmd, checking mana and applying troop effects
func (gs *GameSession) handleDeploy(cmd DeployCmd) {
	mutex.Lock()
	defer mutex.Unlock()
	p := gs.Players[cmd.PlayerIndex]

	spec, ok := gs.TroopSpecs[cmd.TroopName] // stats lookup
	if !ok || p.Mana < spec.Mana {
		return // invalid or insufficient mana
	}
	p.Mana -= spec.Mana
	// apply troop action: attack or heal
	if cmd.TroopName == "Queen" {
		p.HealWeakestTower(300)
	} else {
		gs.attackOpponentTower(cmd.PlayerIndex, spec)
	}
}

// attackOpponentTower resolves a troop attacking the next enemy tower
func (gs *GameSession) attackOpponentTower(idx int, spec TroopSpec) {
	target := gs.Players[1-idx].NextAliveTower()

	// Apply level multiplier to attack
	baseATK := float64(spec.ATK) * gs.Players[idx].Level.Multiplier

	// Calculate critical hit
	isCrit := rand.Float64() < 0.1 // 10% crit chance
	if isCrit {
		baseATK *= 1.2 // 20% more damage on crit
	}

	// Calculate final damage
	dmg := int(baseATK) - target.DEF
	if dmg < 0 {
		dmg = 0
	}

	target.HP -= dmg
	if target.HP <= 0 {
		gs.Players[1-idx].DestroyTower(target)
		gs.justDestroyedTower = true

		// Award EXP for tower destruction
		gs.awardExp(idx, target)
	} else {
		gs.justDestroyedTower = false
	}
}

// awardExp handles EXP gain and leveling
func (gs *GameSession) awardExp(playerIdx int, tower *Tower) {
	player := gs.Players[playerIdx]

	// Award EXP based on tower type
	expGain := 0
	switch tower.Name {
	case "King Tower":
		expGain = 200
	case "Guard Tower":
		expGain = 100
	}

	// Add EXP and check for level up
	player.Level.Exp += expGain
	gs.checkLevelUp(player)
}

// checkLevelUp handles player level progression
func (gs *GameSession) checkLevelUp(player *Player) {
	for player.Level.Exp >= player.Level.NextLevel {
		player.Level.Level++
		player.Level.Exp -= player.Level.NextLevel
		player.Level.NextLevel = int(float64(player.Level.NextLevel) * 1.1) // 10% increase
		player.Level.Multiplier = 1.0 + (float64(player.Level.Level) * 0.1) // 10% per level

		// Notify client of level up
		if player.Conn != nil {
			levelData := fmt.Sprintf(`{"level":%d,"exp":%d,"next_level":%d,"multiplier":%.2f}`,
				player.Level.Level, player.Level.Exp, player.Level.NextLevel, player.Level.Multiplier)
			SendPDU(player.Conn, PDU{
				Type: "level_up",
				Data: []byte(levelData),
			})
		}
	}
}

// broadcastState would serialize and send STATE_UPDATE to clients
func (gs *GameSession) broadcastState() {
	state := GameState{
		YourMana:       gs.Players[0].Mana,
		OpponentMana:   gs.Players[1].Mana,
		YourTowers:     make([]Tower, 0),
		OpponentTowers: make([]Tower, 0),
	}

	// Add player 0's towers
	for _, t := range gs.Players[0].Towers {
		if t.HP > 0 {
			state.YourTowers = append(state.YourTowers, Tower{
				Name: t.Name,
				HP:   t.HP,
			})
		}
	}

	// Add player 1's towers
	for _, t := range gs.Players[1].Towers {
		if t.HP > 0 {
			state.OpponentTowers = append(state.OpponentTowers, Tower{
				Name: t.Name,
				HP:   t.HP,
			})
		}
	}

	// Serialize and send state update
	data, err := json.Marshal(state)
	if err != nil {
		return
	}

	// Send to both players
	for _, p := range gs.Players {
		if p.Conn != nil {
			SendPDU(p.Conn, PDU{
				Type: "state_update",
				Data: data,
			})
		}
	}
}

// checkGameEnd returns true if a King Tower is destroyed
func (gs *GameSession) checkGameEnd() bool {
	for _, p := range gs.Players {
		if p.KingTowerDestroyed() {
			return true
		}
	}
	return false
}

// evaluateWinner compares towers on timeout and assigns EXP
func (gs *GameSession) evaluateWinner() {
	mutex.Lock()
	defer mutex.Unlock()

	// Count remaining towers for each player
	towers0 := 0
	towers1 := 0
	for _, t := range gs.Players[0].Towers {
		if t.HP > 0 {
			towers0++
		}
	}
	for _, t := range gs.Players[1].Towers {
		if t.HP > 0 {
			towers1++
		}
	}

	// Determine winner and assign EXP
	var winner, loser *Player
	if towers0 > towers1 {
		winner = gs.Players[0]
		loser = gs.Players[1]
	} else if towers1 > towers0 {
		winner = gs.Players[1]
		loser = gs.Players[0]
	} else {
		// Draw - both get small EXP
		for _, p := range gs.Players {
			if p.Conn != nil {
				SendPDU(p.Conn, PDU{
					Type: "game_end",
					Data: json.RawMessage(`{"result":"draw","exp":10}`),
				})
			}
		}
		return
	}

	// Winner gets more EXP
	if winner.Conn != nil {
		SendPDU(winner.Conn, PDU{
			Type: "game_end",
			Data: json.RawMessage(`{"result":"win","exp":30}`),
		})
	}
	if loser.Conn != nil {
		SendPDU(loser.Conn, PDU{
			Type: "game_end",
			Data: json.RawMessage(`{"result":"loss","exp":5}`),
		})
	}
}

// NewGameSession creates a new game session
func NewGameSession(mode GameMode, players [2]*Player,
	troopSpecs map[string]TroopSpec,
	towerSpecs map[string]TowerSpec) *GameSession {

	// Initialize player levels
	for i := range players {
		players[i].Level = Level{
			Level:      1,
			Exp:        0,
			NextLevel:  100,
			Multiplier: 1.0,
		}
	}

	return &GameSession{
		Mode:         mode,
		Players:      players,
		TroopSpecs:   troopSpecs,
		TowerSpecs:   towerSpecs,
		Commands:     make(chan DeployCmd, 100),
		Done:         make(chan struct{}),
		TickInterval: time.Second,
	}
}
