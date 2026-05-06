package db

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"task-cli/ent"
	"task-cli/ent/enttest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func setupTestRepo(t *testing.T) (*Repository, *ent.Client) {
	// Using sqlite3 as the driver name as it is what Ent expects for modernc.org/sqlite
	client := enttest.Open(t, "sqlite3", "file:testdb?mode=memory&cache=shared&_fk=1")
	return NewRepository(client), client
}

func TestCreateTask_Success(t *testing.T) {
	repo, client := setupTestRepo(t)
	defer client.Close()

	ctx := context.Background()
	title := "Test Task"
	task, err := repo.CreateTask(ctx, title)

	require.NoError(t, err)
	assert.NotZero(t, task.ID)
	assert.Equal(t, title, task.Title)
	assert.Equal(t, "todo", task.Status)
	assert.NotZero(t, task.CreatedAt)
	assert.NotZero(t, task.UpdatedAt)
}

func TestCreateTask_EmptyTitle(t *testing.T) {
	repo, client := setupTestRepo(t)
	defer client.Close()

	ctx := context.Background()
	_, err := repo.CreateTask(ctx, "")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation")
}

func TestCreateTask_TitleTooLong(t *testing.T) {
	repo, client := setupTestRepo(t)
	defer client.Close()

	ctx := context.Background()
	title := ""
	for i := 0; i < 101; i++ {
		title += "a"
	}
	_, err := repo.CreateTask(ctx, title)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation")
}

func TestGetTaskByID_NotFound(t *testing.T) {
	repo, client := setupTestRepo(t)
	defer client.Close()

	ctx := context.Background()
	_, err := repo.GetTaskByID(ctx, 999)

	assert.ErrorIs(t, err, ErrNotFound)
}

func TestUpdateStatus_Valid(t *testing.T) {
	repo, client := setupTestRepo(t)
	defer client.Close()

	ctx := context.Background()
	task, err := repo.CreateTask(ctx, "Test Task")
	require.NoError(t, err)

	time.Sleep(10 * time.Millisecond) // Ensure updated_at is different

	updated, err := repo.UpdateStatus(ctx, task.ID, "doing")
	require.NoError(t, err)
	assert.Equal(t, "doing", updated.Status)
	assert.True(t, updated.UpdatedAt.After(task.UpdatedAt))
}

func TestUpdateStatus_InvalidEnum(t *testing.T) {
	repo, client := setupTestRepo(t)
	defer client.Close()

	ctx := context.Background()
	task, err := repo.CreateTask(ctx, "Test Task")
	require.NoError(t, err)

	_, err = repo.UpdateStatus(ctx, task.ID, "invalid")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation")
}

func TestGetAllTasks_Order(t *testing.T) {
	repo, client := setupTestRepo(t)
	defer client.Close()

	ctx := context.Background()
	_, err := repo.CreateTask(ctx, "Task 1")
	require.NoError(t, err)
	time.Sleep(10 * time.Millisecond)
	_, err = repo.CreateTask(ctx, "Task 2")
	require.NoError(t, err)

	tasks, err := repo.GetAllTasks(ctx)
	require.NoError(t, err)
	require.Len(t, tasks, 2)
	assert.Equal(t, "Task 2", tasks[0].Title) // Descending order
	assert.Equal(t, "Task 1", tasks[1].Title)
}

func TestEnsureInitialized(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gask-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	t.Run("Missing", func(t *testing.T) {
		err := EnsureInitialized(tmpDir)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Run 'gask init' first")
	})

	t.Run("Present", func(t *testing.T) {
		gaskDir := filepath.Join(tmpDir, ".gask")
		err := os.Mkdir(gaskDir, 0755)
		require.NoError(t, err)

		dbFile := filepath.Join(gaskDir, "gask.db")
		err = os.WriteFile(dbFile, []byte("fake db"), 0644)
		require.NoError(t, err)

		err = EnsureInitialized(tmpDir)
		assert.NoError(t, err)
	})
}
