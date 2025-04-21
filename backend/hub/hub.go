package hub

import (
	"log"
	"sync"

	"file-converter/backend/types"
)

type Hub struct {
	Clients    map[*types.Client]bool
	Broadcast  chan types.Message
	Register   chan *types.Client
	Unregister chan *types.Client
	mutex      sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[*types.Client]bool),
		Broadcast:  make(chan types.Message),
		Register:   make(chan *types.Client),
		Unregister: make(chan *types.Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mutex.Lock()
			h.Clients[client] = true
			h.mutex.Unlock()

			// Broadcast user joined message
			h.Broadcast <- types.Message{
				Type: types.UserPresence,
				Payload: map[string]string{
					"type":     "join",
					"username": client.Username,
				},
			}

		case client := <-h.Unregister:
			h.mutex.Lock()
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				client.Conn.Close()
			}
			h.mutex.Unlock()

			// Broadcast user left message
			h.Broadcast <- types.Message{
				Type: types.UserPresence,
				Payload: map[string]string{
					"type":     "leave",
					"username": client.Username,
				},
			}

		case message := <-h.Broadcast:
			h.mutex.Lock()
			for client := range h.Clients {
				err := client.Conn.WriteJSON(message)
				if err != nil {
					log.Printf("error: %v", err)
					client.Conn.Close()
					delete(h.Clients, client)
				}
			}
			h.mutex.Unlock()
		}
	}
}
