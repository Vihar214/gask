package cmd

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddCommand_ValidateTitle(t *testing.T) {
	tempDir := t.TempDir()
	cwd, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(cwd)

	os.Mkdir(".gask", 0755)

	tests := []struct {
		name    string
		args    []string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "empty title",
			args:    []string{""},
			wantErr: true,
			errMsg:  "title is required",
		},
		{
			name:    "too long title",
			args:    []string{"this title is definitely more than one hundred characters long and should be rejected by the validation logic that we implemented in the command."},
			wantErr: true,
			errMsg:  "title cannot exceed 100 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := addCmd.RunE(addCmd, tt.args)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
