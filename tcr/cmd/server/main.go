package main

import (
	"tcr/server"
)

func main() {
	// Load users
	users, err := server.LoadUsers("players.json")
	if err != nil {
		panic("failed to load users: " + err.Error())
	}
	// Load specs
	troops := make(map[string]server.TroopSpec)
	towers := make(map[string]server.TowerSpec)
	if err := server.LoadSpecs("specs.json", troops, towers); err != nil {
		panic("failed to load specs: " + err.Error())
	}

	matchQueue := make(chan *server.ClientHandler, 10) // adjust size as needed

	// Start the server
	if err := server.StartServer(":9000", users, troops, towers, matchQueue); err != nil {
		panic("failed to start server: " + err.Error())
	}

}
