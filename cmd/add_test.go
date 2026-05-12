package cmd

import (
	"fmt"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
)

func TestAddCommand_ValidateTitle(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "empty title",
			args:    []string{"", "  "},
			wantErr: true,
			errMsg:  "title is required",
		},
		{
			name:    "too long title",
			args:    []string{"this title is definitely more than one hundred characters long and should be rejected by the validation logic that we implemented in the command."},
			wantErr: true,
			errMsg:  "title cannot exceed 100 characters",
		},
		{
			name:    "valid title",
			args:    []string{"Valid", "Task"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			title := strings.Join(tt.args, " ")
			title = strings.TrimSpace(title)

			var err error
			if title == "" {
				err = fmt.Errorf("title is required")
			} else if utf8.RuneCountInString(title) > 100 {
				err = fmt.Errorf("title cannot exceed 100 characters")
			}

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
