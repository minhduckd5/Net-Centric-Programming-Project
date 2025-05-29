package main

import (
	"tcr/server"
	"tcr/specs"
)

func main() {
	// Load users
	users, err := server.LoadUsers("players.json")
	if err != nil {
		panic("failed to load users: " + err.Error())
	}
	// log.Println(users)
	// Load specs
	loadedSpecs, err := specs.LoadSpecs("../../specs/game_specs.json")
	if err != nil {
		panic("failed to load specs: " + err.Error())
	}

	troops := loadedSpecs.Troops
	towers := loadedSpecs.Towers
	// log.Println("troops: ", troops, "towers: ", towers)
	matchQueue := make(chan *server.ClientHandler, 10) // adjust size as needed

	// Start the server
	if err := server.StartServer(":9000", users, troops, towers, matchQueue); err != nil {
		panic("failed to start server: " + err.Error())
	}

}
