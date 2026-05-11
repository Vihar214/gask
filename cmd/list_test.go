package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTruncateTitle(t *testing.T) {
	// Exact length
	title := "This title is exactly fifty characters long right." // 50 chars
	assert.Equal(t, title, truncateTitle(title, 50))

	// Truncation
	longTitle := "This title is much longer than fifty characters and should be truncated."
	truncated := truncateTitle(longTitle, 50)
	assert.True(t, len(truncated) <= 50)
	assert.Contains(t, truncated, "…")

	// Empty
	assert.Equal(t, "", truncateTitle("", 50))
}

func TestTruncateDescription(t *testing.T) {
	// Empty
	assert.Equal(t, "—", truncateDescription("", 80))

	// Truncation
	longDesc := "This is a very long description that definitely exceeds eighty characters and must be cut down."
	truncated := truncateDescription(longDesc, 80)
	assert.True(t, len(truncated) <= 80)
	assert.Contains(t, truncated, "…")
}

func TestFormatStatus(t *testing.T) {
	// formatStatus returns the string wrapped in lipgloss styling, 
	// which doesn't include the color name in the output string.
	// We verify that the returned string is not empty and is a valid status string.
	assert.NotEmpty(t, formatStatus("todo"))
	assert.Contains(t, formatStatus("todo"), "todo")
	assert.Contains(t, formatStatus("doing"), "doing")
	assert.Contains(t, formatStatus("done"), "done")
	assert.Contains(t, formatStatus("blocked"), "blocked")
	assert.Contains(t, formatStatus("invalid"), "invalid")
}
