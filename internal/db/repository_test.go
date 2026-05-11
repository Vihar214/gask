package db

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"task-cli/ent"
	"task-cli/ent/enttest"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestRepo(t *testing.T) (*Repository, *ent.Client) {
	client := enttest.Open(t, "sqlite3", "file:testdb?mode=memory&cache=shared&_fk=1")
	return NewRepository(client), client
}

func createTaskWithStatus(t *testing.T, repo *Repository, ctx context.Context, title, status string) *Task {
	t.Helper()
	task, err := repo.CreateTask(ctx, title)
	if err != nil {
		t.Fatalf("failed to create task: %v", err)
	}
	updated, err := repo.UpdateStatus(ctx, task.ID, status)
	if err != nil {
		t.Fatalf("failed to update task status: %v", err)
	}
	return updated
}

func TestCreateTask_Success(t *testing.T) {
	repo, client := setupTestRepo(t)
	defer func() {
		require.NoError(t, client.Close())
	}()

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
	defer func() {
		require.NoError(t, client.Close())
	}()

	ctx := context.Background()
	_, err := repo.CreateTask(ctx, "")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation")
}

func TestCreateTask_TitleTooLong(t *testing.T) {
	repo, client := setupTestRepo(t)
	defer func() {
		require.NoError(t, client.Close())
	}()

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
	defer func() {
		require.NoError(t, client.Close())
	}()

	ctx := context.Background()
	_, err := repo.GetTaskByID(ctx, 999)

	assert.ErrorIs(t, err, ErrNotFound)
}

func TestListTasks_SortOrder(t *testing.T) {
	repo, client := setupTestRepo(t)
	defer func() {
		require.NoError(t, client.Close())
	}()

	ctx := context.Background()

	// Create tasks in random order with different statuses
	createTaskWithStatus(t, repo, ctx, "task 1 (todo)", "todo")
	createTaskWithStatus(t, repo, ctx, "task 2 (doing)", "doing")
	createTaskWithStatus(t, repo, ctx, "task 3 (done)", "done")
	createTaskWithStatus(t, repo, ctx, "task 4 (blocked)", "blocked")

	tasks, err := repo.ListTasks(ctx)
	require.NoError(t, err)

	assert.Equal(t, 4, len(tasks))
	// Expected order: doing → todo → blocked → done
	assert.Equal(t, "doing", tasks[0].Status)
	assert.Equal(t, "todo", tasks[1].Status)
	assert.Equal(t, "blocked", tasks[2].Status)
	assert.Equal(t, "done", tasks[3].Status)
}

func TestListTasks_Empty(t *testing.T) {
	repo, client := setupTestRepo(t)
	defer func() {
		require.NoError(t, client.Close())
	}()

	ctx := context.Background()
	tasks, err := repo.ListTasks(ctx)
	require.NoError(t, err)
	assert.Empty(t, tasks)
}

func TestUpdateStatus_InvalidEnum(t *testing.T) {
	repo, client := setupTestRepo(t)
	defer func() {
		require.NoError(t, client.Close())
	}()

	ctx := context.Background()
	task, err := repo.CreateTask(ctx, "Test Task")
	require.NoError(t, err)

	_, err = repo.UpdateStatus(ctx, task.ID, "invalid")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation")
}

func TestUpdateDescription_TooLong(t *testing.T) {
	repo, client := setupTestRepo(t)
	defer func() {
		require.NoError(t, client.Close())
	}()

	ctx := context.Background()
	task, err := repo.CreateTask(ctx, "Test Task")
	require.NoError(t, err)

	description := ""
	for i := 0; i < 10001; i++ {
		description += "a"
	}
	_, err = repo.UpdateDescription(ctx, task.ID, description)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation")
}

func TestEnsureInitialized(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gask-test-*")
	require.NoError(t, err)
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Fatalf("failed to remove tmpDir %q: %v", tmpDir, err)
		}
	}()

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
