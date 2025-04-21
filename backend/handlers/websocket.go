package handlers

import (
	"log"
	"net/http"

	"file-converter/backend/hub"
	"file-converter/backend/types"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // For development only
	},
}

type WebSocketHandler struct {
	hub *hub.Hub
}

func NewWebSocketHandler(h *hub.Hub) *WebSocketHandler {
	return &WebSocketHandler{hub: h}
}

func (h *WebSocketHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &types.Client{
		ID:   uuid.New().String(),
		Conn: conn,
	}

	h.hub.Register <- client
	go h.handleMessages(client)
}

func (h *WebSocketHandler) handleMessages(client *types.Client) {
	defer func() {
		h.hub.Unregister <- client
	}()

	for {
		var message types.Message
		err := client.Conn.ReadJSON(&message)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		// Handle different message types
		switch message.Type {
		case types.ChatMessage:
			if payload, ok := message.Payload.(map[string]interface{}); ok {
				if username, exists := payload["username"].(string); exists && client.Username == "" {
					client.Username = username
				}
			}
			h.hub.Broadcast <- message

		case types.GameState:
			// Handle game state updates
			h.hub.Broadcast <- message

		case types.UserPresence:
			// Handle user presence updates
			h.hub.Broadcast <- message
		}
	}
}
