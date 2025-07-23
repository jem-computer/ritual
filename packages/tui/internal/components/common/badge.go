// ABOUTME: Badge component for displaying status indicators with colored backgrounds
// ABOUTME: Used for task status (ACTIVE/PAUSED) and log status (SUCCESS/ERROR)

package common

import (
	"github.com/jem-computer/ritual/tui/internal/styles"
	"github.com/jem-computer/ritual/tui/internal/theme"
)

// BadgeType represents the type of badge
type BadgeType int

const (
	BadgeActive BadgeType = iota
	BadgePaused
	BadgeSuccess
	BadgeError
	BadgeWarning
	BadgeInfo
)

// Badge renders a status badge with appropriate styling
func Badge(text string, badgeType BadgeType) string {
	t := theme.CurrentTheme()
	if t == nil {
		return text
	}

	style := styles.NewStyle().
		Padding(0, 1).
		Bold(true)

	switch badgeType {
	case BadgeActive:
		style = style.Background(t.Success()).Foreground(t.Background())
	case BadgePaused:
		style = style.Background(t.Warning()).Foreground(t.Background())
	case BadgeSuccess:
		style = style.Background(t.Success()).Foreground(t.Background())
	case BadgeError:
		style = style.Background(t.Error()).Foreground(t.Background())
	case BadgeWarning:
		style = style.Background(t.Warning()).Foreground(t.Background())
	case BadgeInfo:
		style = style.Background(t.Info()).Foreground(t.Background())
	}

	return style.Render(text)
}

// StatusBadge renders a task status badge
func StatusBadge(status string) string {
	switch status {
	case "ACTIVE":
		return Badge(status, BadgeActive)
	case "PAUSED":
		return Badge(status, BadgePaused)
	default:
		return Badge(status, BadgeInfo)
	}
}

// LogStatusBadge renders a log status badge
func LogStatusBadge(status string) string {
	switch status {
	case "SUCCESS":
		return Badge(status, BadgeSuccess)
	case "ERROR":
		return Badge(status, BadgeError)
	default:
		return Badge(status, BadgeInfo)
	}
}
