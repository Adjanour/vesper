package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	_ "modernc.org/sqlite"
)

const (
	dataDir       = "./data"
	dbFile        = "./data/tasks.db"
	migrationsDir = "./internal/database/migrations"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run migrate.go [up|down]")
	}

	command := os.Args[1]

	// Ensure data directory exists
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}

	// Connect to database
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	switch command {
	case "up":
		if err := migrateUp(db); err != nil {
			log.Fatalf("Migration up failed: %v", err)
		}
		log.Println("✓ Migrations applied successfully")
	case "down":
		if err := migrateDown(db); err != nil {
			log.Fatalf("Migration down failed: %v", err)
		}
		log.Println("✓ Migrations rolled back successfully")
	default:
		log.Fatalf("Unknown command: %s. Use 'up' or 'down'", command)
	}
}

func migrateUp(db *sql.DB) error {
	files, err := getMigrationFiles(".up.sql")
	if err != nil {
		return err
	}

	for _, file := range files {
		log.Printf("Applying migration: %s", file)
		content, err := os.ReadFile(filepath.Join(migrationsDir, file))
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", file, err)
		}

		if _, err := db.Exec(string(content)); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", file, err)
		}
	}

	return nil
}

func migrateDown(db *sql.DB) error {
	files, err := getMigrationFiles(".down.sql")
	if err != nil {
		return err
	}

	// Reverse order for down migrations
	sort.Sort(sort.Reverse(sort.StringSlice(files)))

	for _, file := range files {
		log.Printf("Rolling back migration: %s", file)
		content, err := os.ReadFile(filepath.Join(migrationsDir, file))
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", file, err)
		}

		if _, err := db.Exec(string(content)); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", file, err)
		}
	}

	return nil
}

func getMigrationFiles(suffix string) ([]string, error) {
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read migrations directory: %w", err)
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), suffix) {
			files = append(files, entry.Name())
		}
	}

	sort.Strings(files)
	return files, nil
}
