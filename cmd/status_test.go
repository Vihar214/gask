package cmd

import (
	"context"
	"os"
	"path/filepath"
	"task-cli/internal/db"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStatusCmd(t *testing.T) {
	tempDir := t.TempDir()
	os.MkdirAll(filepath.Join(tempDir, ".gask"), 0755)
	dbPath := filepath.Join(tempDir, ".gask", "gask.db")
	cwd, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(cwd)

	client, err := db.Open(context.Background(), dbPath)
	require.NoError(t, err)
	defer client.Close()
	repo := db.NewRepository(client)

	task, err := repo.CreateTask(context.Background(), "test task")
	require.NoError(t, err)

	t.Run("Valid update", func(t *testing.T) {
		err := statusCmd.RunE(statusCmd, []string{string(rune(task.ID + '0')), "doing"})
		assert.NoError(t, err)
		updatedTask, err := repo.GetTaskByID(context.Background(), task.ID)
		assert.NoError(t, err)
		assert.Equal(t, "doing", updatedTask.Status)
	})

	t.Run("Invalid status", func(t *testing.T) {
		err := statusCmd.RunE(statusCmd, []string{string(rune(task.ID + '0')), "invalid"})
		assert.Error(t, err)
	})

	t.Run("Done to Todo guard", func(t *testing.T) {
		_, err := repo.UpdateStatus(context.Background(), task.ID, "done")
		require.NoError(t, err)
		err = statusCmd.RunE(statusCmd, []string{string(rune(task.ID + '0')), "todo"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "without --force")
	})

	t.Run("Done to Todo with force", func(t *testing.T) {
		forceStatus = true
		defer func() { forceStatus = false }()
		err := statusCmd.RunE(statusCmd, []string{string(rune(task.ID + '0')), "todo"})
		assert.NoError(t, err)
		updatedTask, err := repo.GetTaskByID(context.Background(), task.ID)
		assert.NoError(t, err)
		assert.Equal(t, "todo", updatedTask.Status)
	})
}
