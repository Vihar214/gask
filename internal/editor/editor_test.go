package editor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
	assert.Contains(t, err.Error(), "VS Code")
	assert.Contains(t, err.Error(), "Install 'code' command in PATH")

	// Test empty input
	path, err = Validate("")
	assert.Error(t, err)
	assert.Equal(t, "editor command cannot be empty", err.Error())
}
