package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"task-cli/internal/db"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var long bool

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tasks",
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		dbPath := filepath.Join(cwd, ".gask", "gask.db")
		client, err := db.Open(context.Background(), dbPath)
		if err != nil {
			return err
		}
		defer client.Close()
		repo := db.NewRepository(client)

		ctx := context.Background()
		tasks, err := repo.ListTasks(ctx)
		if err != nil {
			return err
		}

		if len(tasks) == 0 {
			fmt.Println("No tasks found in this directory. Run 'gask add' to create one.")
			return nil
		}

		renderTable(tasks, long)
		return nil
	},
}

func init() {
	listCmd.Flags().BoolVarP(&long, "long", "l", false, "Show description")
	rootCmd.AddCommand(listCmd)
}

func renderTable(tasks []*db.Task, showDescription bool) {
	width := 100 // default
	if w, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
		width = w
	}

	headers := []string{"ID", "Status", "Title"}
	if showDescription {
		headers = append(headers, "Description")
	}
	headers = append(headers, "Age")

	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
		Headers(headers...).
		Width(width).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == -1 { // Header row
				return lipgloss.NewStyle().Bold(true).Underline(true)
			}
			if col == 0 { // ID column
				return lipgloss.NewStyle().Align(lipgloss.Right)
			}
			return lipgloss.NewStyle()
		})

	for _, task := range tasks {
		row := []string{
			fmt.Sprintf("%d", task.ID),
			formatStatus(task.Status),
			truncateTitle(task.Title, 50),
		}
		if showDescription {
			row = append(row, truncateDescription(task.Description, 80))
		}
		row = append(row, FormatAge(task.CreatedAt))
		t.Row(row...)
	}

	fmt.Println(t.Render())
}

func formatStatus(status string) string {
	colors := map[string]string{
		"todo":    "blue",
		"doing":   "yellow",
		"done":    "green",
		"blocked": "red",
	}
	color, ok := colors[status]
	if !ok {
		color = "white"
	}
	return lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Render(status)
}

func truncateTitle(title string, maxLen int) string {
	if len(title) <= maxLen {
		return title
	}
	runes := []rune(title)
	truncated := string(runes[:maxLen-1])
	lastSpace := strings.LastIndex(truncated, " ")
	if lastSpace != -1 {
		truncated = strings.TrimSpace(truncated[:lastSpace])
	}
	return truncated + "…"
}

func truncateDescription(desc string, maxLen int) string {
	if desc == "" {
		return "—"
	}
	if len(desc) <= maxLen {
		return desc
	}
	runes := []rune(desc)
	truncated := string(runes[:maxLen-1])
	lastSpace := strings.LastIndex(truncated, " ")
	if lastSpace != -1 {
		truncated = strings.TrimSpace(truncated[:lastSpace])
	}
	return truncated + "…"
}
