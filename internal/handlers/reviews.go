package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/dmltdev/flashcards/internal/models"
)

func (h *Handler) CreateReview(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	cardID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Error("Invalid card ID", err)
		http.Error(w, "Invalid card ID", http.StatusBadRequest)
		return
	}

	var review models.Review
	if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
		log.Error("Invalid JSON", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	review.CardID = cardID
	review.ReviewedAt = time.Now()

	if err := review.Validate(); err != nil {
		log.Error("Invalid review", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if review.Quality >= 4 {
		review.NextReviewAt = time.Now().Add(3 * 24 * time.Hour)
	} else {
		review.NextReviewAt = time.Now().Add(24 * time.Hour)
	}

	if err := h.db.CreateReview(&review); err != nil {
		log.Error("Failed to create review", err)
		http.Error(w, "Failed to create review", http.StatusInternalServerError)
		return
	}

	log.Info("Review created", "review", review)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(review)
}
