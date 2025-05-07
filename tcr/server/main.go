// server/main.go
package main

import (
	"flag"
	"log"
)

func main() {
	port := flag.String("port", "9000", "TCP port")
	users := flag.String("users", "players.json", "User data JSON")
	flag.Parse()

	addr := ":" + *port
	if err := StartServer(addr, *users); err != nil {
		log.Fatal(err)
	}
}
