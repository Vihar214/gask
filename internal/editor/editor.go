package editor

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// ErrNoEditor is returned when no editor can be resolved from environment or config.
var ErrNoEditor = errors.New("no editor configured; set $EDITOR or run 'gask init'")

// EditorRunner is a function that returns an *exec.Cmd.
// It is used for mocking the editor execution in tests.
var EditorRunner = exec.Command

// ResolveEditor returns the editor to use, prioritizing $EDITOR over the provided configEditor.
func ResolveEditor(configEditor string) (string, error) {
	if envEditor := os.Getenv("EDITOR"); envEditor != "" {
		return envEditor, nil
	}
	if configEditor != "" {
		return configEditor, nil
	}
	return "", ErrNoEditor
}

// Open spawns the editor to open the specified file.
// It pipes Stdin, Stdout, and Stderr to the current process.
func Open(editorPath string, filePath string) error {
	cmd := EditorRunner(editorPath, filePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("editor failed: %w", err)
	}
	return nil
}

// Validate checks if the given editor command is available in the system PATH.
// It returns the full path to the executable if found, or an error with helpful hints.
func Validate(name string) (string, error) {
	if name == "" {
		return "", fmt.Errorf("editor command cannot be empty")
	}

	// Simple case: just the command name
	path, err := exec.LookPath(name)
	if err == nil {
		return path, nil
	}

	// If not found, provide hints for common editors
	if strings.Contains(strings.ToLower(name), "code") {
		return "", fmt.Errorf("editor '%s' not found. If you are using VS Code, ensure the 'code' command is in your PATH. (In VS Code, open Command Palette and run 'Shell Command: Install 'code' command in PATH')", name)
	}

	return "", fmt.Errorf("editor '%s' not found in PATH", name)
}
