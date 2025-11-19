package database

import (
	"database/sql"
	"fmt"

	"github.com/dmltdev/flashcards/internal/models"
)

func (db *DB) CreateCard(card *models.Card) error {
	query := `
		INSERT INTO cards (deck_id, front, back, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id, created_at, updated_at`

	err := db.QueryRow(query, card.DeckID, card.Front, card.Back).Scan(
		&card.ID, &card.CreatedAt, &card.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create card: %w", err)
	}
	return nil
}

func (db *DB) GetCard(id int) (*models.Card, error) {
	var card models.Card
	query := `SELECT id, deck_id, front, back, created_at, updated_at FROM cards WHERE id = $1`
	
	err := db.Get(&card, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("card not found")
		}
		return nil, fmt.Errorf("failed to get card: %w", err)
	}
	return &card, nil
}

func (db *DB) GetCardsByDeck(deckID int) ([]models.Card, error) {
	var cards []models.Card
	query := `SELECT id, deck_id, front, back, created_at, updated_at FROM cards WHERE deck_id = $1 ORDER BY created_at DESC`
	
	err := db.Select(&cards, query, deckID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cards by deck: %w", err)
	}
	return cards, nil
}

func (db *DB) UpdateCard(card *models.Card) error {
	query := `
		UPDATE cards 
		SET front = $1, back = $2, updated_at = NOW()
		WHERE id = $3
		RETURNING updated_at`

	err := db.QueryRow(query, card.Front, card.Back, card.ID).Scan(&card.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("card not found")
		}
		return fmt.Errorf("failed to update card: %w", err)
	}
	return nil
}

func (db *DB) DeleteCard(id int) error {
	query := `DELETE FROM cards WHERE id = $1`
	
	result, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete card: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("card not found")
	}

	return nil
}
