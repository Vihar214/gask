package db

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpen_SchemaFailure(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gask-test-*")
	require.NoError(t, err)
	t.Cleanup(func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Logf("failed to remove tmpDir: %v", err)
		}
	})

	// Create a file where the DB should be, to cause migration failure
	dbPath := filepath.Join(tmpDir, "gask.db")
	err = os.WriteFile(dbPath, []byte("not a sqlite db"), 0644)
	require.NoError(t, err)

	// Make the file read-only to ensure schema creation fails
	err = os.Chmod(dbPath, 0444)
	require.NoError(t, err)

	ctx := context.Background()
	_, err = Open(ctx, dbPath)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed creating schema resources")
}
