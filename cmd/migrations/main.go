package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dmltdev/flashcards/internal/database"
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

	if len(os.Args) < 2 {
		fmt.Println("Usage: go run cmd/migrations/main.go [up|down]")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "up":
		if err := runMigrationsUp(db); err != nil {
			log.Fatal("Failed to run migrations up:", err)
		}
		fmt.Println("Migrations completed successfully!")
	case "down":
		if err := runMigrationsDown(db); err != nil {
			log.Fatal("Failed to run migrations down:", err)
		}
		fmt.Println("Migrations rolled back successfully!")
	default:
		fmt.Println("Invalid command. Use 'up' or 'down'")
		os.Exit(1)
	}
}

func runMigrationsUp(db *database.DB) error {
	if err := createMigrationsTable(db); err != nil {
		return err
	}

	migrations := []Migration{
		{Name: "001_create_decks_table", Up: createDecksTable},
		{Name: "002_create_cards_table", Up: createCardsTable},
		{Name: "003_create_reviews_table", Up: createReviewsTable},
	}

	for _, migration := range migrations {
		if err := runMigration(db, migration); err != nil {
			return fmt.Errorf("failed to run migration %s: %w", migration.Name, err)
		}
	}

	return nil
}

func runMigrationsDown(db *database.DB) error {
	// Drop tables in reverse order
	queries := []string{
		"DROP TABLE IF EXISTS reviews CASCADE;",
		"DROP TABLE IF EXISTS cards CASCADE;",
		"DROP TABLE IF EXISTS decks CASCADE;",
		"DROP TABLE IF EXISTS migrations CASCADE;",
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("failed to execute: %s, error: %w", query, err)
		}
	}

	return nil
}

type Migration struct {
	Name string
	Up   func(*database.DB) error
}

func createMigrationsTable(db *database.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS migrations (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL UNIQUE,
			executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`

	_, err := db.Exec(query)
	return err
}

func runMigration(db *database.DB, migration Migration) error {
	// Check if migration already ran
	var count int
	err := db.Get(&count, "SELECT COUNT(*) FROM migrations WHERE name = $1", migration.Name)
	if err != nil {
		return err
	}

	if count > 0 {
		fmt.Printf("Migration %s already executed, skipping...\n", migration.Name)
		return nil
	}

	if err := migration.Up(db); err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO migrations (name) VALUES ($1)", migration.Name)
	if err != nil {
		return err
	}

	fmt.Printf("Migration %s executed successfully\n", migration.Name)
	return nil
}

func createDecksTable(db *database.DB) error {
	query := `
		CREATE TABLE decks (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TRIGGER update_decks_updated_at
			BEFORE UPDATE ON decks
			FOR EACH ROW
			EXECUTE FUNCTION update_updated_at_column();`

	_, err := db.Exec(query)
	return err
}

func createCardsTable(db *database.DB) error {
	query := `
		CREATE TABLE cards (
			id SERIAL PRIMARY KEY,
			deck_id INTEGER NOT NULL REFERENCES decks(id) ON DELETE CASCADE,
			front TEXT NOT NULL,
			back TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);

		CREATE INDEX idx_cards_deck_id ON cards(deck_id);

		CREATE TRIGGER update_cards_updated_at
			BEFORE UPDATE ON cards
			FOR EACH ROW
			EXECUTE FUNCTION update_updated_at_column();`

	_, err := db.Exec(query)
	return err
}

func createReviewsTable(db *database.DB) error {
	query := `
		CREATE TABLE reviews (
			id SERIAL PRIMARY KEY,
			card_id INTEGER NOT NULL REFERENCES cards(id) ON DELETE CASCADE,
			quality INTEGER NOT NULL CHECK (quality >= 1 AND quality <= 5),
			reviewed_at TIMESTAMP NOT NULL,
			next_review_at TIMESTAMP NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);

		CREATE INDEX idx_reviews_card_id ON reviews(card_id);
		CREATE INDEX idx_reviews_next_review_at ON reviews(next_review_at);

		CREATE TRIGGER update_reviews_updated_at
			BEFORE UPDATE ON reviews
			FOR EACH ROW
			EXECUTE FUNCTION update_updated_at_column();`

	_, err := db.Exec(query)
	return err
}