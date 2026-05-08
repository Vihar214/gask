package cmd

import (
	"os"

	"task-cli/internal/db"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gask",
	Short: "gask is a minimalist task manager for developers",
	Long:  `gask is a minimalist task manager for developers that lives in your project's .gask directory.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Skip for 'init' command and 'help'
		if cmd.Name() == "init" || cmd.Name() == "help" {
			return nil
		}
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		return db.EnsureInitialized(cwd)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Root flags can be added here
}
