package main

import (
	"log"
	"net/http"
	"os"

	"file-converter/backend/handlers"

	"file-converter/backend/hub"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Create new hub
	h := hub.NewHub()
	go h.Run()

	// Create WebSocket handler
	wsHandler := handlers.NewWebSocketHandler(h)

	// Set up routes
	http.HandleFunc("/ws", wsHandler.HandleWebSocket)

	// Set up CORS middleware
	corsMiddleware := func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			handler.ServeHTTP(w, r)
		})
	}

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	addr := ":" + port
	log.Printf("Server starting on port %s", port)

	server := &http.Server{
		Addr:    addr,
		Handler: corsMiddleware(http.DefaultServeMux),
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
