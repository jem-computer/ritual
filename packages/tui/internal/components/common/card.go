// ABOUTME: Card component for displaying task information with status badges
// ABOUTME: Renders task cards with metadata rows and action buttons

package common

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss/v2"
	"github.com/jem-computer/ritual/tui/internal/styles"
	"github.com/jem-computer/ritual/tui/internal/theme"
)

// Task represents a scheduled task
type Task struct {
	ID       string
	Name     string
	Status   string // ACTIVE, PAUSED
	Prompt   string
	Schedule string
	Output   string
	Model    string
	NextRun  string
	LastRun  string
}

// TaskCard renders a task as a styled card
func TaskCard(task Task, width int, selected bool) string {
	t := theme.CurrentTheme()
	if t == nil {
		return "No theme"
	}

	// Card container style
	cardStyle := styles.NewStyle().
		Width(width).
		Padding(1).
		MarginBottom(1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.BorderSubtle())

	if selected {
		cardStyle = cardStyle.BorderForeground(t.BorderActive())
	}

	// Build card content
	var content strings.Builder

	// Status badge and task name on same line
	statusBadge := StatusBadge(task.Status)
	taskName := styles.NewStyle().
		Foreground(t.Text()).
		Bold(true).
		Render(task.Name)

	// Action buttons (placeholder for now)
	actions := styles.NewStyle().
		Foreground(t.TextMuted()).
		Render("⏸ ✏ 🗑")

	// Calculate available width for task name
	statusWidth := lipgloss.Width(statusBadge)
	actionsWidth := lipgloss.Width(actions)
	availableWidth := width - statusWidth - actionsWidth - 6 // padding and spaces

	// Header line with status, name, and actions
	header := fmt.Sprintf("%s  %s%s%s",
		statusBadge,
		taskName,
		strings.Repeat(" ", max(1, availableWidth-lipgloss.Width(taskName))),
		actions,
	)
	content.WriteString(header)
	content.WriteString("\n\n")

	// Metadata rows
	labelStyle := styles.NewStyle().
		Foreground(t.TextMuted()).
		Width(10)

	valueStyle := styles.NewStyle().
		Foreground(t.Text())

	// PROMPT row
	content.WriteString(labelStyle.Render("PROMPT:"))
	content.WriteString(" ")
	content.WriteString(valueStyle.Render(truncate(task.Prompt, width-12)))
	content.WriteString("\n")

	// SCHEDULE row
	content.WriteString(labelStyle.Render("SCHEDULE:"))
	content.WriteString(" ")
	content.WriteString(valueStyle.Render(fmt.Sprintf("%s | OUTPUT: %s | MODEL: %s",
		task.Schedule, task.Output, task.Model)))
	content.WriteString("\n")

	// Run times row
	runStyle := styles.NewStyle().
		Foreground(t.TextMuted()).
		Faint(true)

	content.WriteString(runStyle.Render(fmt.Sprintf("NEXT RUN: %s | LAST RUN: %s",
		task.NextRun, task.LastRun)))

	return cardStyle.Render(content.String())
}

// truncate truncates a string to a maximum length with ellipsis
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return "..."
	}
	return s[:maxLen-3] + "..."
}

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
