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
			return fmt.Errorf("invalid task ID: '%s'", args[0])
		}
		newStatus := args[1]
		validStatuses := map[string]bool{"todo": true, "doing": true, "done": true, "blocked": true}
		if !validStatuses[newStatus] {
			return fmt.Errorf("invalid status: '%s'", newStatus)
		}

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
		defer func() {
			if err := client.Close(); err != nil {
				fmt.Fprintf(os.Stderr, "failed to close database: %v\n", err)
			}
		}()

		repo := db.NewRepository(client)

		task, err := repo.GetTaskByID(ctx, id)
		if err != nil {
			return err
		}

		if task.Status == "done" && newStatus == "todo" && !forceStatus {
			return fmt.Errorf("task #%d is already done. Are you sure you want to move it back to todo?\nRun: gask status %d todo --force", id, id)
		}

		_, err = repo.UpdateStatus(ctx, id, newStatus)
		if err != nil {
			return err
		}

		fmt.Printf("✓ Task #%d status updated: %s → %s\n", id, task.Status, newStatus)
		return nil
	},
}

func init() {
	statusCmd.Flags().BoolVarP(&forceStatus, "force", "f", false, "Force status transition")
	rootCmd.AddCommand(statusCmd)
}
