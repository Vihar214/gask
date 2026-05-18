package editor

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolveEditor(t *testing.T) {
	// Backup and restore EDITOR env var
	oldEditor := os.Getenv("EDITOR")
	defer func() {
		if err := os.Setenv("EDITOR", oldEditor); err != nil {
			t.Fatalf("failed to restore EDITOR: %v", err)
		}
	}()

	t.Run("prioritize $EDITOR", func(t *testing.T) {
		if err := os.Setenv("EDITOR", "vim"); err != nil {
			t.Fatalf("failed to set EDITOR: %v", err)
		}
		editor, err := ResolveEditor("nano")
		assert.NoError(t, err)
		assert.Equal(t, "vim", editor)
	})

	t.Run("fallback to config", func(t *testing.T) {
		if err := os.Setenv("EDITOR", ""); err != nil {
			t.Fatalf("failed to set EDITOR: %v", err)
		}
		editor, err := ResolveEditor("nano")
		assert.NoError(t, err)
		assert.Equal(t, "nano", editor)
	})

	t.Run("return error if both empty", func(t *testing.T) {
		if err := os.Setenv("EDITOR", ""); err != nil {
			t.Fatalf("failed to set EDITOR: %v", err)
		}
		editor, err := ResolveEditor("")
		assert.ErrorIs(t, err, ErrNoEditor)
		assert.Empty(t, editor)
	})
}

func TestOpen(t *testing.T) {
	// Mock EditorRunner
	oldRunner := EditorRunner
	defer func() { EditorRunner = oldRunner }()

	EditorRunner = func(name string, arg ...string) *exec.Cmd {
		// Return a command that succeeds immediately
		return exec.Command("true")
	}

	err := Open("vim", "test.md")
	assert.NoError(t, err)

	EditorRunner = func(name string, arg ...string) *exec.Cmd {
		// Return a command that fails
		return exec.Command("false")
	}

	err = Open("vim", "test.md")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "editor failed")
}

func TestValidate(t *testing.T) {
	// Test common editors (assuming they might or might not be in PATH)
	// We can't guarantee 'vim' or 'nano' is in PATH on all test environments,
	// but 'ls' or 'sh' usually is.
	path, err := Validate("sh")
	assert.NoError(t, err)
	assert.NotEmpty(t, path)

	// Test non-existent command
	path, err = Validate("non-existent-command-12345")
	assert.Error(t, err)
	assert.Empty(t, path)
	assert.Contains(t, err.Error(), "not found in PATH")

	// Test VS Code hint
	path, err = Validate("code-not-here")
	assert.Error(t, err)
	assert.Empty(t, path)
	assert.Contains(t, err.Error(), "VS Code")
	assert.Contains(t, err.Error(), "Install 'code' command in PATH")

	// Test empty input
	path, err = Validate("")
	assert.Error(t, err)
	assert.Empty(t, path)
	assert.Equal(t, "editor command cannot be empty", err.Error())
}
