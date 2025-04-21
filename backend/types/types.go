package types

import "github.com/gorilla/websocket"

type MessageType string

const (
	ChatMessage  MessageType = "chat"
	GameState    MessageType = "game_state"
	UserPresence MessageType = "user_presence"
)

type Message struct {
	Type    MessageType `json:"type"`
	Payload interface{} `json:"payload"`
}

type ChatPayload struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

type Client struct {
	ID       string
	Username string
	Conn     *websocket.Conn
}
