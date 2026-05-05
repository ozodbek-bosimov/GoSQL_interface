package database

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
)

func RunMigrations() error {
	db := GetDB()
	if db == nil {
		return fmt.Errorf("database not initialized")
	}

	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	migrations, err := getMigrationFiles()
	if err != nil {
		return fmt.Errorf("failed to get migration files: %w", err)
	}

	for _, migration := range migrations {
		var exists bool
		err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = $1)", migration).Scan(&exists)
		if err != nil {
			return fmt.Errorf("failed to check migration status: %w", err)
		}

		if exists {
			log.Printf("  Migration %s already applied, skipping", migration)
			continue
		}

		content, err := os.ReadFile(filepath.Join("migrations", migration))
		if err != nil {
			return fmt.Errorf("failed to read migration %s: %w", migration, err)
		}

		_, err = db.Exec(string(content))
		if err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", migration, err)
		}

		_, err = db.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", migration)
		if err != nil {
			return fmt.Errorf("failed to record migration %s: %w", migration, err)
		}

		log.Printf("✓ Applied migration: %s", migration)
	}

	log.Println("✓ All migrations completed")
	return nil
}

func getMigrationFiles() ([]string, error) {
	files, err := os.ReadDir("migrations")
	if err != nil {
		return nil, err
	}

	var migrations []string
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".sql" {
			migrations = append(migrations, file.Name())
		}
	}

	sort.Strings(migrations)
	return migrations, nil
}
