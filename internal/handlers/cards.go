package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/dmltdev/flashcards/internal/database"
	"github.com/dmltdev/flashcards/internal/models"
)

type Handler struct {
	db *database.DB
}

func NewHandler(db *database.DB) *Handler {
	return &Handler{db: db}
}

func (h *Handler) CreateDeck(w http.ResponseWriter, r *http.Request) {
	var deck models.Deck
	if err := json.NewDecoder(r.Body).Decode(&deck); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := deck.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.db.CreateDeck(&deck); err != nil {
		http.Error(w, "Failed to create deck", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(deck)
}

func (h *Handler) GetDeck(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid deck ID", http.StatusBadRequest)
		return
	}

	deck, err := h.db.GetDeck(id)
	if err != nil {
		http.Error(w, "Deck not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deck)
}

func (h *Handler) CreateCard(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	deckID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid deck ID", http.StatusBadRequest)
		return
	}

	var card models.Card
	if err := json.NewDecoder(r.Body).Decode(&card); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	card.DeckID = deckID

	if err := card.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.db.CreateCard(&card); err != nil {
		http.Error(w, "Failed to create card", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(card)
}

func (h *Handler) GetNextCard(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	deckID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid deck ID", http.StatusBadRequest)
		return
	}

	cards, err := h.db.GetCardsByDeck(deckID)
	if err != nil {
		http.Error(w, "Failed to get cards", http.StatusInternalServerError)
		return
	}

	if len(cards) == 0 {
		http.Error(w, "No cards found in deck", http.StatusNotFound)
		return
	}

	// For MVP: return the first card (oldest created)
	// Later: implement spaced repetition logic
	nextCard := cards[len(cards)-1] // Get oldest card

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nextCard)
}
