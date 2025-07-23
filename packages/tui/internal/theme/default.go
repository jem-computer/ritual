// ABOUTME: Default theme implementation for Ritual with OpenCode-inspired colors
// ABOUTME: Provides a dark theme with navy backgrounds and vibrant accent colors

package theme

import (
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/charmbracelet/lipgloss/v2/compat"
)

// DefaultTheme is the default theme for Ritual
type DefaultTheme struct {
	BaseTheme
	name string
}

// NewDefaultTheme creates a new instance of the default theme
func NewDefaultTheme() *DefaultTheme {
	theme := &DefaultTheme{
		name: "ritual",
	}

	// Background colors - dark navy palette
	theme.BackgroundColor = compat.AdaptiveColor{
		Dark:  lipgloss.Color("#0a0a0a"),
		Light: lipgloss.Color("#ffffff"),
	}
	theme.BackgroundPanelColor = compat.AdaptiveColor{
		Dark:  lipgloss.Color("#141414"),
		Light: lipgloss.Color("#fafafa"),
	}
	theme.BackgroundElementColor = compat.AdaptiveColor{
		Dark:  lipgloss.Color("#1e1e1e"),
		Light: lipgloss.Color("#f5f5f5"),
	}

	// Border colors
	theme.BorderSubtleColor = compat.AdaptiveColor{
		Dark:  lipgloss.Color("#3c3c3c"),
		Light: lipgloss.Color("#d4d4d4"),
	}
	theme.BorderColor = compat.AdaptiveColor{
		Dark:  lipgloss.Color("#484848"),
		Light: lipgloss.Color("#b8b8b8"),
	}
	theme.BorderActiveColor = compat.AdaptiveColor{
		Dark:  lipgloss.Color("#606060"),
		Light: lipgloss.Color("#a0a0a0"),
	}

	// Brand colors
	theme.PrimaryColor = compat.AdaptiveColor{
		Dark:  lipgloss.Color("#fab283"), // Warm orange
		Light: lipgloss.Color("#3b7dd8"),
	}
	theme.SecondaryColor = compat.AdaptiveColor{
		Dark:  lipgloss.Color("#5c9cf5"), // Blue
		Light: lipgloss.Color("#7b5bb6"),
	}
	theme.AccentColor = compat.AdaptiveColor{
		Dark:  lipgloss.Color("#9d7cd8"), // Purple
		Light: lipgloss.Color("#d68c27"),
	}

	// Text colors
	theme.TextColor = compat.AdaptiveColor{
		Dark:  lipgloss.Color("#eeeeee"),
		Light: lipgloss.Color("#1a1a1a"),
	}
	theme.TextMutedColor = compat.AdaptiveColor{
		Dark:  lipgloss.Color("#808080"),
		Light: lipgloss.Color("#8a8a8a"),
	}

	// Status colors
	theme.ErrorColor = compat.AdaptiveColor{
		Dark:  lipgloss.Color("#e06c75"), // Red
		Light: lipgloss.Color("#d1383d"),
	}
	theme.WarningColor = compat.AdaptiveColor{
		Dark:  lipgloss.Color("#f5a742"), // Orange
		Light: lipgloss.Color("#d68c27"),
	}
	theme.SuccessColor = compat.AdaptiveColor{
		Dark:  lipgloss.Color("#7fd88f"), // Green
		Light: lipgloss.Color("#3d9a57"),
	}
	theme.InfoColor = compat.AdaptiveColor{
		Dark:  lipgloss.Color("#56b6c2"), // Cyan
		Light: lipgloss.Color("#318795"),
	}

	return theme
}

// Name returns the name of the theme
func (t *DefaultTheme) Name() string {
	return t.name
}

func init() {
	// Register the default theme
	RegisterTheme("ritual", NewDefaultTheme())
}
