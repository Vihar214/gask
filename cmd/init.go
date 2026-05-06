package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"task-cli/internal/db"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialise gask in the current directory",
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("init: failed to get current directory: %w", err)
		}

		// 1. Check if .gask.db (legacy file) exists
		legacyPath := filepath.Join(cwd, ".gask.db")
		if _, err := os.Stat(legacyPath); err == nil {
			return fmt.Errorf("Found legacy .gask.db file. Please move it to .gask/gask.db or delete it before running init.")
		}

		// 2. Check if .gask/ already exists
		gaskDir := filepath.Join(cwd, ".gask")
		if _, err := os.Stat(gaskDir); err == nil {
			fmt.Println("gask is already initialised in this directory.")
			return nil
		}

		// 3. Create .gask/ directory
		if err := os.Mkdir(gaskDir, 0755); err != nil {
			return fmt.Errorf("init: cannot create .gask/: %w", err)
		}

		// 4. Open the DB connection (triggers auto-migration)
		dbPath := filepath.Join(gaskDir, "gask.db")
		client, err := db.Open(context.Background(), dbPath)
		if err != nil {
			return fmt.Errorf("init: migration failed: %w", err)
		}
		defer client.Close()

		// 5. Print success message and the .gitignore hint
		fmt.Printf("✓ Initialised gask in .gask/\n")
		fmt.Printf("  DB: .gask/gask.db\n\n")
		fmt.Printf("Hint: add .gask/ to your .gitignore to keep tasks out of version control.\n")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
