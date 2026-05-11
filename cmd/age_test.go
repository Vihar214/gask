package cmd

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFormatAge(t *testing.T) {
	now := time.Now()
	assert.Equal(t, "just now", FormatAge(now.Add(-30*time.Second)))
	assert.Equal(t, "5m ago", FormatAge(now.Add(-5*time.Minute)))
	assert.Equal(t, "3h ago", FormatAge(now.Add(-3*time.Hour)))
	assert.Equal(t, "2d ago", FormatAge(now.Add(-48*time.Hour)))
	assert.Equal(t, "1w ago", FormatAge(now.Add(-13*24*time.Hour)))
	assert.Equal(t, "2w ago", FormatAge(now.Add(-14*24*time.Hour)))
}
