package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig_SaveAndLoad(t *testing.T) {
	// Create a temporary directory for the test
	tmpDir, err := os.MkdirTemp("", "gask-config-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Change working directory to the temp dir
	oldCwd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(oldCwd)
	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	// Test Load when config doesn't exist
	cfg, err := Load()
	assert.ErrorIs(t, err, ErrConfigNotFound)
	assert.Nil(t, cfg)

	// Test Save
	cfg = &Config{Editor: "vim"}
	err = cfg.Save()
	assert.NoError(t, err)

	// Verify file exists and has correct content
	configPath := filepath.Join(tmpDir, ".gask", "config.json")
	assert.FileExists(t, configPath)

	data, err := os.ReadFile(configPath)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"editor": "vim"`)

	// Test Load
	loadedCfg, err := Load()
	assert.NoError(t, err)
	assert.Equal(t, cfg.Editor, loadedCfg.Editor)
}
