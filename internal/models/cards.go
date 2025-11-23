package models

import (
	"errors"
	"strings"
	"time"
)

type Card struct {
	ID        int       `json:"id" db:"id"`
    DeckID    int       `json:"deck_id" db:"deck_id"`
    Front     string    `json:"front" db:"front"`
    Back      string    `json:"back" db:"back"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Deck struct {
    ID int `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	Cards []Card `json:"cards,omitempty" db:"-"`
	CardCount int `json:"card_count" db:"card_count"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Review struct {
	ID int `json:"id" db:"id"`
	CardID int `json:"card_id" db:"card_id"`
	Quality int `json:"quality" db:"quality"`
	ReviewedAt time.Time `json:"reviewed_at" db:"reviewed_at"`
	NextReviewAt time.Time `json:"next_review_at" db:"next_review_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func (c *Card) Validate() error {
	if strings.TrimSpace(c.Front) == "" {
		return errors.New("front cannot be empty")
	}
	if strings.TrimSpace(c.Back) == "" {
		return errors.New("back cannot be empty")
	}
	if c.DeckID <= 0 {
		return errors.New("deck_id must be positive")
	}
	return nil
}

func (d *Deck) Validate() error {
	if strings.TrimSpace(d.Name) == "" {
		return errors.New("name cannot be empty")
	}
	return nil
}

func (r *Review) Validate() error {
	if r.Quality < 1 || r.Quality > 5 {
		return errors.New("quality must be between 1 and 5")
	}
	if r.CardID <= 0 {
		return errors.New("card_id must be positive")
	}
	return nil
}