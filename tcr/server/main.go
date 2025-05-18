// server/main.go
package main

import (
	// "tcr/common"
	"flag"
	"log"
)

func main() {
	port := flag.String("port", "9000", "TCP port")
	usersFile := flag.String("users", "players.json", "User data JSON")
	specsFile := flag.String("specs", "specs.json", "Unit specs JSON")
	flag.Parse()

	// 1. Load users
	userData, err := common.LoadUsers(*usersFile)
	if err != nil {
		log.Fatalf("Failed to load users: %v", err)
	}

	// 2. Load specs
	troops := make(map[string]common.TroopSpec)
	towers := make(map[string]common.TowerSpec)
	if err := LoadSpecs(*specsFile, troops, towers); err != nil {
		log.Fatalf("Failed to load specs: %v", err)
	}

	// 3. Initialize matchmaking queue
	matchQueue := make(chan *ClientHandler)

	// 4. Start server with matchmaking
	addr := ":" + *port
	if err := StartServer(addr, userData, troops, towers, matchQueue); err != nil {
		log.Fatal(err)
	}
}
