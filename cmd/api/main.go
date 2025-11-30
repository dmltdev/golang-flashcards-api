package main

import (
	"log"
	"net/http"
	"os"

	"github.com/dmltdev/flashcards/internal/database"
	"github.com/dmltdev/flashcards/internal/handlers"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	db, err := database.NewConnection()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	handler := handlers.NewHandler(db)

	mux := http.NewServeMux()
	
	mux.HandleFunc("POST /decks", handler.CreateDeck)

	mux.HandleFunc("GET /decks", handler.GetDecks)
	mux.HandleFunc("GET /decks/{id}", handler.GetDeck)
	
	mux.HandleFunc("POST /decks/{id}/cards", handler.CreateCard)
	mux.HandleFunc("GET /decks/{id}/cards/next", handler.GetNextCard)
	mux.HandleFunc("POST /cards/{id}/reviews", handler.CreateReview)

	port := getEnv("SERVER_PORT", "8080")
	log.Printf("Server starting on port %s", port)
	
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}