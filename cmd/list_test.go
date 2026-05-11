package cmd

import (
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
)

func TestTruncateTitle(t *testing.T) {
	// Exact length
	title := "This title is exactly fifty characters long right." // 50 chars
	assert.Equal(t, title, truncateTitle(title, 50))

	// Truncation
	longTitle := "This title is much longer than fifty characters and should be truncated."
	truncated := truncateTitle(longTitle, 50)
	assert.True(t, utf8.RuneCountInString(truncated) <= 50)
	assert.Contains(t, truncated, "…")

	// Unicode
	unicodeTitle := "🚀🚀🚀"
	assert.Equal(t, unicodeTitle, truncateTitle(unicodeTitle, 5))
	assert.Equal(t, "🚀…", truncateTitle(unicodeTitle, 2))

	// Guards
	assert.Equal(t, "…", truncateTitle("anything", 0))
	assert.Equal(t, "…", truncateTitle("anything", -1))

	// Empty
	assert.Equal(t, "", truncateTitle("", 50))
}

func TestTruncateDescription(t *testing.T) {
	// Empty
	assert.Equal(t, "—", truncateDescription("", 80))

	// Truncation
	longDesc := "This is a very long description that definitely exceeds eighty characters and must be cut down."
	truncated := truncateDescription(longDesc, 80)
	assert.True(t, utf8.RuneCountInString(truncated) <= 80)
	assert.Contains(t, truncated, "…")

	// Unicode
	unicodeDesc := "🚀🚀🚀"
	assert.Equal(t, unicodeDesc, truncateDescription(unicodeDesc, 5))
	assert.Equal(t, "🚀…", truncateDescription(unicodeDesc, 2))

	// Guards
	assert.Equal(t, "…", truncateDescription("anything", 0))
	assert.Equal(t, "…", truncateDescription("anything", -1))
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
