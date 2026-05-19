package cmd

import (
	"context"
	"os"
	"path/filepath"
	"strconv"
	"task-cli/internal/db"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStatusCmd(t *testing.T) {
	tempDir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(tempDir, ".gask"), 0755); err != nil {
		t.Fatalf("failed to create directory: %v", err)
	}
	dbPath := filepath.Join(tempDir, ".gask", "gask.db")
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}
	defer func() {
		if err := os.Chdir(cwd); err != nil {
			t.Errorf("failed to restore working directory: %v", err)
		}
	}()

	client, err := db.Open(context.Background(), dbPath)
	require.NoError(t, err)
	defer func() {
		if err := client.Close(); err != nil {
			t.Errorf("failed to close database client: %v", err)
		}
	}()
	repo := db.NewRepository(client)

	task, err := repo.CreateTask(context.Background(), "test task")
	require.NoError(t, err)

	t.Run("Valid update", func(t *testing.T) {
		err := statusCmd.RunE(statusCmd, []string{strconv.Itoa(task.ID), "doing"})
		assert.NoError(t, err)
		updatedTask, err := repo.GetTaskByID(context.Background(), task.ID)
		assert.NoError(t, err)
		assert.Equal(t, "doing", updatedTask.Status)
	})

	t.Run("Invalid status", func(t *testing.T) {
		err := statusCmd.RunE(statusCmd, []string{strconv.Itoa(task.ID), "invalid"})
		assert.Error(t, err)
	})

	t.Run("Done to Todo guard", func(t *testing.T) {
		_, err := repo.UpdateStatus(context.Background(), task.ID, "done")
		require.NoError(t, err)
		err = statusCmd.RunE(statusCmd, []string{strconv.Itoa(task.ID), "todo"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Are you sure you want to move it back to todo?")
	})

	t.Run("Done to Todo with force", func(t *testing.T) {
		forceStatus = true
		defer func() { forceStatus = false }()
		err := statusCmd.RunE(statusCmd, []string{strconv.Itoa(task.ID), "todo"})
		assert.NoError(t, err)
		updatedTask, err := repo.GetTaskByID(context.Background(), task.ID)
		assert.NoError(t, err)
		assert.Equal(t, "todo", updatedTask.Status)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		err := statusCmd.RunE(statusCmd, []string{"not-an-id", "todo"})
		assert.Error(t, err)
	})

	t.Run("Non-existent ID", func(t *testing.T) {
		err := statusCmd.RunE(statusCmd, []string{strconv.Itoa(task.ID + 999), "doing"})
		assert.Error(t, err)
	})

	t.Run("Done to Doing", func(t *testing.T) {
		_, err := repo.UpdateStatus(context.Background(), task.ID, "done")
		require.NoError(t, err)
		err = statusCmd.RunE(statusCmd, []string{strconv.Itoa(task.ID), "doing"})
		assert.NoError(t, err)
		updatedTask, err := repo.GetTaskByID(context.Background(), task.ID)
		assert.NoError(t, err)
		assert.Equal(t, "doing", updatedTask.Status)
	})

	t.Run("Done to Blocked", func(t *testing.T) {
		_, err := repo.UpdateStatus(context.Background(), task.ID, "done")
		require.NoError(t, err)
		err = statusCmd.RunE(statusCmd, []string{strconv.Itoa(task.ID), "blocked"})
		assert.NoError(t, err)
		updatedTask, err := repo.GetTaskByID(context.Background(), task.ID)
		assert.NoError(t, err)
		assert.Equal(t, "blocked", updatedTask.Status)
	})
}
