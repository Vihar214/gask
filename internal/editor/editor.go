package editor

import (
	"fmt"
	"os/exec"
	"strings"
)

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
