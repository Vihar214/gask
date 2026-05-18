package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"task-cli/internal/db"

	"github.com/spf13/cobra"
)

var forceStatus bool

var statusCmd = &cobra.Command{
	Use:   "status [id] [status]",
	Short: "Update the status of a task",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid task ID: %w", err)
		}
		newStatus := args[1]

		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		dbPath := filepath.Join(cwd, ".gask", "gask.db")
		ctx := cmd.Context()
		if ctx == nil {
			ctx = context.Background()
		}
		client, err := db.Open(ctx, dbPath)
		if err != nil {
			return err
		}
		defer client.Close()

		repo := db.NewRepository(client)

		task, err := repo.GetTaskByID(ctx, id)
		if err != nil {
			return err
		}

		if task.Status == "done" && newStatus == "todo" && !forceStatus {
			return fmt.Errorf("cannot transition from 'done' to 'todo' without --force")
		}

		_, err = repo.UpdateStatus(ctx, id, newStatus)
		if err != nil {
			return err
		}

		fmt.Printf("Task %d updated to '%s'\n", id, newStatus)
		return nil
	},
}

func init() {
	statusCmd.Flags().BoolVarP(&forceStatus, "force", "f", false, "Force status transition")
	rootCmd.AddCommand(statusCmd)
}
