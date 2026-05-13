package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"task-cli/internal/config"
	"task-cli/internal/db"
	"task-cli/internal/editor"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add [title]",
	Short: "Add a new task",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if _, err := os.Stat(".gask"); os.IsNotExist(err) {
			return fmt.Errorf("Error: gask not initialized. Run 'gask init' first")
		}

		title := strings.Join(args, " ")
		title = strings.TrimSpace(title)

		if title == "" {
			return fmt.Errorf("title is required")
		}
		if utf8.RuneCountInString(title) > 100 {
			return fmt.Errorf("title cannot exceed 100 characters")
		}

		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		editorPath, err := editor.ResolveEditor(cfg.Editor)
		if err != nil {
			return err
		}

		// Create temp file for description
		tmpFile, err := os.CreateTemp("", "gask-*.md")
		if err != nil {
			return fmt.Errorf("failed to create temp file: %w", err)
		}
		tmpPath := tmpFile.Name()
		if err := tmpFile.Close(); err != nil {
			return fmt.Errorf("failed to close temp file: %w", err)
		}
		defer os.Remove(tmpPath)

		// Open editor
		if err := editor.Open(editorPath, tmpPath); err != nil {
			return err
		}

		// Read description back
		descBytes, err := os.ReadFile(tmpPath)
		if err != nil {
			return fmt.Errorf("failed to read description: %w", err)
		}
		description := strings.TrimSpace(string(descBytes))

		// Save to DB
		dbPath := filepath.Join(cwd, ".gask", "gask.db")
		client, err := db.Open(context.Background(), dbPath)
		if err != nil {
			return err
		}
		defer func() {
			closeErr := client.Close()
			if closeErr != nil {
				if err == nil {
					err = closeErr
				} else {
					err = fmt.Errorf("%w; close error: %v", err, closeErr)
				}
			}
		}()
		repo := db.NewRepository(client)

		ctx := context.Background()
		task, err := repo.CreateTask(ctx, title)
		if err != nil {
			return fmt.Errorf("failed to create task: %w", err)
		}

		if description != "" {
			_, err = repo.UpdateDescription(ctx, task.ID, description)
			if err != nil {
				return fmt.Errorf("failed to update task description: %w", err)
			}
		}

		fmt.Printf("✓ Task #%d created: \"%s\"\n", task.ID, title)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
