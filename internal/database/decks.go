package database

import (
	"database/sql"
	"fmt"

	"github.com/dmltdev/flashcards/internal/models"
)

func (db *DB) CreateDeck(deck *models.Deck) error {
	query := `
		INSERT INTO decks (name, created_at, updated_at)
		VALUES ($1, NOW(), NOW())
		RETURNING id, created_at, updated_at`

	err := db.QueryRow(query, deck.Name).Scan(
		&deck.ID, &deck.CreatedAt, &deck.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create deck: %w", err)
	}
	return nil
}

func (db *DB) GetDeck(id int) (*models.Deck, error) {
	var deck models.Deck
	query := `SELECT id, name, created_at, updated_at FROM decks WHERE id = $1`
	
	err := db.Get(&deck, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("deck not found")
		}
		return nil, fmt.Errorf("failed to get deck: %w", err)
	}
	return &deck, nil
}

func (db *DB) GetAllDecks() ([]models.Deck, error) {
	var decks []models.Deck
	query := `SELECT id, name, created_at, updated_at FROM decks ORDER BY created_at DESC`
	
	err := db.Select(&decks, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get decks: %w", err)
	}
	return decks, nil
}

func (db *DB) UpdateDeck(deck *models.Deck) error {
	query := `
		UPDATE decks 
		SET name = $1, updated_at = NOW()
		WHERE id = $2
		RETURNING updated_at`

	err := db.QueryRow(query, deck.Name, deck.ID).Scan(&deck.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("deck not found")
		}
		return fmt.Errorf("failed to update deck: %w", err)
	}
	return nil
}

func (db *DB) DeleteDeck(id int) error {
	query := `DELETE FROM decks WHERE id = $1`
	
	result, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete deck: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("deck not found")
	}

	return nil
}