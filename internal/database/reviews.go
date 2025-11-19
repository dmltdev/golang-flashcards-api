package database

import (
	"fmt"

	"github.com/dmltdev/flashcards/internal/models"
)

func (db *DB) CreateReview(review *models.Review) error {
	query := `
		INSERT INTO reviews (card_id, quality, reviewed_at, next_review_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id, created_at, updated_at`

	err := db.QueryRow(query, review.CardID, review.Quality, review.ReviewedAt, review.NextReviewAt).Scan(
		&review.ID, &review.CreatedAt, &review.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create review: %w", err)
	}
	return nil
}

func (db *DB) GetReviewsByCard(cardID int) ([]models.Review, error) {
	var reviews []models.Review
	query := `SELECT id, card_id, quality, reviewed_at, next_review_at, created_at, updated_at 
			  FROM reviews WHERE card_id = $1 ORDER BY reviewed_at DESC`
	
	err := db.Select(&reviews, query, cardID)
	if err != nil {
		return nil, fmt.Errorf("failed to get reviews by card: %w", err)
	}
	return reviews, nil
}