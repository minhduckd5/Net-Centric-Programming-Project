// game.go
// Core game logic for Text-Based Clash Royale (TCR)

package server

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"sync"
	"tcr/specs"
	"time"
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
	Players            [2]*Player // two players
	TroopSpecs         map[string]specs.TroopSpec
	TowerSpecs         map[string]specs.TowerSpec
	Commands           chan DeployCmd // incoming deploy commands
	Done               chan struct{}  // signals end of game
	TickInterval       time.Duration  // for enhanced mode
	justDestroyedTower bool           // tracks if a tower was just destroyed
}

type TroopInstance struct {
	Spec   specs.TroopSpec
	Health int
	// Possibly: Position, OwnerIndex, SpawnTime, etc.
}

// DeployCmd is issued by a client or AI to deploy a troop
type DeployCmd struct {
	PlayerIndex int    // 0 or 1
	TroopName   string // e.g., "Pawn"
}

// GameState represents the current state of the game
type GameState struct {
	YourMana      int               `json:"your_mana"`
	OpponentMana  int               `json:"opponent_mana"`
	Player1Towers []specs.TowerSpec `json:"your_towers"`
	Player2Towers []specs.TowerSpec `json:"opponent_towers"`
}

// startGame launches the appropriate game loop based on mode
func (gs *GameSession) StartGame() {
	for i, player := range gs.Players {
		go func(index int, conn net.Conn) {
			for {
				pdu, err := ReceivePDU(conn)
				if err != nil {
					log.Printf("Error receiving PDU: %v", err)
					return
				}

				switch pdu.Type {
				case "deploy":
					var payload struct {
						Troop string `json:"troop"`
					}
					if err := json.Unmarshal(pdu.Data, &payload); err != nil {
						log.Println("Invalid deploy payload:", err)
						continue
					}

					gs.Commands <- DeployCmd{
						PlayerIndex: index,
						TroopName:   payload.Troop,
					}
				}
			}
		}(i, player.Conn)
	}

	// Start game loop
	gs.enhancedLoop()

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

			// 🔽 Check if game has ended after tick
			if gs.checkGameEnd() {
				gs.evaluateWinner()
				close(gs.Done)
				return
			}

		case cmd := <-gs.Commands:
			gs.handleDeploy(cmd)

			// 🔽 Check if game has ended after deploy
			if gs.checkGameEnd() {
				gs.evaluateWinner()
				close(gs.Done)
				return
			}

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

	for i, p := range gs.Players {
		if p.Mana < 10 {
			p.Mana++
		}

		// Each tower attacks one troop (if any)
		for _, tower := range p.Towers {
			opponent := gs.Players[1-i]
			if len(opponent.ActiveTroops) == 0 {
				continue
			}
			target := opponent.ActiveTroops[0]
			// Apply level multiplier to attack

			// Basic target selection: first troop
			// Apply level multiplier to attack
			baseATK := float64(tower.Damage) * gs.Players[i].Level.Multiplier

			// Calculate critical hit
			isCrit := rand.Float64() < 0.1 // 10% crit chance
			if isCrit {
				baseATK *= 1.2 // 20% more damage on crit
			}

			// Calculate final damage
			dmg := max(int(baseATK)-target.Spec.Defence, 0)
			target.Health -= dmg
			if target.Health <= 0 {
				opponent.ActiveTroops = opponent.ActiveTroops[1:]
				log.Println("Troop die, active list: ", opponent.ActiveTroops)
			}
		}
	}
	gs.broadcastState()
}

// handleDeploy processes a DeployCmd, checking mana and applying troop effects
func (gs *GameSession) handleDeploy(cmd DeployCmd) {
	mutex.Lock()
	defer mutex.Unlock()
	//Take the player
	p := gs.Players[cmd.PlayerIndex]

	spec, ok := gs.TroopSpecs[cmd.TroopName] // stats lookup
	// log.Println("Spec to deploy: ", spec)
	if !ok || p.Mana < spec.Cost {
		if !ok {
			log.Println("Troop name: ", cmd.TroopName)
			log.Println("gs.Troop name: ", gs.TroopSpecs)
			log.Println("spec not find")
		} else {
			log.Println("Mana insufficient")
		}
		return // invalid or insufficient mana
	}
	p.Mana -= spec.Cost
	log.Println("Current mana: ", p.Mana)

	// apply troop action: attack or heal
	troop := &TroopInstance{
		Spec:   spec,
		Health: spec.Health,
		// optionally Position, etc.
	}
	p.ActiveTroops = append(p.ActiveTroops, troop)

	if cmd.TroopName == "queen" {
		p.HealWeakestTower(300)
	}
	if cmd.TroopName != "queen" { // Queen already handled as instant support
		go func(troop *TroopInstance, playerIdx int) {
			for {
				time.Sleep(2 * time.Second)

				mutex.Lock()
				if troop.Health <= 0 {
					mutex.Unlock()
					log.Printf("Troop %s is dead, stopping attack loop\n", troop.Spec.Name)
					break
				}
				gs.attackOpponentTowerFromTroop(playerIdx, troop)
				mutex.Unlock()
			}
		}(troop, cmd.PlayerIndex)
	}
}

func (gs *GameSession) attackOpponentTowerFromTroop(playerIdx int, troop *TroopInstance) {
	opponent := gs.Players[1-playerIdx]
	player := gs.Players[playerIdx]

	target := opponent.NextAliveTower()
	if target == nil {
		return
	}

	baseATK := float64(troop.Spec.Damage) * player.Level.Multiplier
	if rand.Float64() < 0.1 {
		baseATK *= 1.2
	}
	dmg := max(int(baseATK)-target.Defence, 0)
	target.Health -= dmg

	log.Printf("Troop %s attacked tower %s for %d damage\n", troop.Spec.Name, target.Name, dmg)

	if target.Health <= 0 {
		opponent.DestroyTower(target)
		gs.justDestroyedTower = true
		gs.awardExp(playerIdx, target)
	} else {
		gs.justDestroyedTower = false
	}
}

// awardExp handles EXP gain and leveling
func (gs *GameSession) awardExp(playerIdx int, tower *specs.TowerSpec) {
	player := gs.Players[playerIdx]

	// Award EXP based on tower type
	expGain := 0
	switch tower.Name {
	case "King Tower":
		expGain = 200
	case "Princess Tower":
		expGain = 150
	case "Guard Tower":
		expGain = 100
	case "Cannon":
		expGain = 50
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
		YourMana:      gs.Players[0].Mana,
		OpponentMana:  gs.Players[1].Mana,
		Player1Towers: make([]specs.TowerSpec, 0),
		Player2Towers: make([]specs.TowerSpec, 0),
	}

	// Add player 0's towers
	for _, t := range gs.Players[0].Towers {
		if t.Health > 0 {
			state.Player1Towers = append(state.Player1Towers, specs.TowerSpec{
				Name:        t.Name,
				Type:        t.Type,
				Health:      t.Health,
				Damage:      t.Damage,
				Defence:     t.Defence,
				Range:       t.Range,
				AttackSpeed: t.AttackSpeed,
				Target:      t.Target,
			})
		}
	}

	// Add player 1's towers
	for _, t := range gs.Players[1].Towers {
		if t.Health > 0 {
			state.Player2Towers = append(state.Player2Towers, specs.TowerSpec{
				Name:        t.Name,
				Type:        t.Type,
				Health:      t.Health,
				Damage:      t.Damage,
				Defence:     t.Defence,
				Range:       t.Range,
				AttackSpeed: t.AttackSpeed,
				Target:      t.Target,
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
		if t.Health > 0 {
			towers0++
		}
	}
	for _, t := range gs.Players[1].Towers {
		if t.Health > 0 {
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
func NewGameSession(players [2]*Player,
	troopSpecs map[string]specs.TroopSpec,
	towerSpecs map[string]specs.TowerSpec) *GameSession {
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
		Players:      players,
		TroopSpecs:   troopSpecs,
		TowerSpecs:   towerSpecs,
		Commands:     make(chan DeployCmd, 100),
		Done:         make(chan struct{}),
		TickInterval: time.Second,
	}
}
