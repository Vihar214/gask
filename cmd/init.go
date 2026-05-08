package cmd

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"task-cli/internal/config"
	"task-cli/internal/db"
	"task-cli/internal/editor"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize gask in the current directory",
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("init: failed to get current directory: %w", err)
		}

		// 1. Check if config already exists (Idempotency)
		cfg, err := config.Load()
		if err == nil {
			fmt.Println("gask is already initialized in this directory.")
			return nil
		}
		if !errors.Is(err, config.ErrConfigNotFound) {
			return fmt.Errorf("init: failed to check config: %w", err)
		}

		// 2. Create .gask/ directory if it doesn't exist
		gaskDir := filepath.Join(cwd, ".gask")
		if _, err := os.Stat(gaskDir); os.IsNotExist(err) {
			if err := os.Mkdir(gaskDir, 0755); err != nil {
				return fmt.Errorf("init: cannot create .gask/: %w", err)
			}
		}

		// 3. Open the DB connection (triggers auto-migration)
		dbPath := filepath.Join(gaskDir, "gask.db")
		client, err := db.Open(context.Background(), dbPath)
		if err != nil {
			return fmt.Errorf("init: migration failed: %w", err)
		}
		defer client.Close()

		// 4. Interactive Editor Prompt
		fmt.Println("Select your preferred editor for tasks:")
		fmt.Println("1. Vim (vim)")
		fmt.Println("2. Nvim (nvim)")
		fmt.Println("3. Nano (nano)")
		fmt.Println("4. VS Code (code)")
		fmt.Println("5. Other (custom command)")

		var selectedEditor string
		scanner := bufio.NewScanner(os.Stdin)

		for attempts := 3; attempts > 0; attempts-- {
			fmt.Printf("Choice [1-5] (%d attempts left): ", attempts)
			if !scanner.Scan() {
				return fmt.Errorf("init: failed to read input")
			}
			choice := strings.TrimSpace(scanner.Text())

			switch choice {
			case "1":
				selectedEditor = "vim"
			case "2":
				selectedEditor = "nvim"
			case "3":
				selectedEditor = "nano"
			case "4":
				selectedEditor = "code"
			case "5":
				fmt.Print("Enter custom editor command: ")
				if !scanner.Scan() {
					return fmt.Errorf("init: failed to read input")
				}
				selectedEditor = strings.TrimSpace(scanner.Text())
			default:
				fmt.Printf("Invalid choice: %s.\n", choice)
				continue
			}

			// Validate selected editor
			_, err := editor.Validate(selectedEditor)
			if err == nil {
				break // Success
			}

			fmt.Printf("Error: %v\n", err)
			selectedEditor = "" // Reset for next loop
			if attempts == 1 {
				return fmt.Errorf("init: failed to configure editor after 3 attempts")
			}
		}

		if selectedEditor == "" {
			return fmt.Errorf("init: failed to configure editor")
		}

		cfg = &config.Config{Editor: selectedEditor}
		if err := cfg.Save(); err != nil {
			return fmt.Errorf("init: failed to save config: %w", err)
		}

		// 5. Print success message and the .gitignore hint
		fmt.Printf("\n✓ initialized gask in .gask/\n")
		fmt.Printf("  DB: .gask/gask.db\n")
		fmt.Printf("  Editor: %s\n\n", selectedEditor)
		fmt.Printf("Hint: add .gask/ to your .gitignore to keep tasks out of version control.\n")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
