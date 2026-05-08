package db

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"task-cli/ent"

	_ "github.com/mattn/go-sqlite3"
)

// Open opens a connection to the SQLite database and runs auto-migrations.
func Open(ctx context.Context, dbPath string) (*ent.Client, error) {
	dsn := fmt.Sprintf("file:%s?_fk=1&_journal_mode=WAL&_synchronous=NORMAL&_busy_timeout=5000", dbPath)
	client, err := ent.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed opening connection to sqlite: %w", err)
	}


	// Run the auto migration tool
	if err := client.Schema.Create(ctx); err != nil {
		client.Close()
		return nil, fmt.Errorf("failed creating schema resources: %w", err)
	}

	return client, nil
}

// EnsureInitialized checks if the .gask/gask.db file exists.
func EnsureInitialized(dir string) error {
	dbPath := filepath.Join(dir, ".gask", "gask.db")
	_, err := os.Stat(dbPath)
	if err == nil {
		return nil
	}
	if os.IsNotExist(err) {
		return fmt.Errorf("gask is not initialised in this directory. Run 'gask init' first")
	}
	return fmt.Errorf("failed to check initialization status of %s: %w", dbPath, err)
}
