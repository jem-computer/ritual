// ABOUTME: Theme interface and base implementation for Ritual's UI theming system
// ABOUTME: Defines color contracts for backgrounds, borders, text, and status indicators

package theme

import (
	"github.com/charmbracelet/lipgloss/v2/compat"
)

// Theme defines the interface for all UI themes in the application.
// All colors must be defined as compat.AdaptiveColor to support
// both light and dark terminal backgrounds.
type Theme interface {
	Name() string

	// Background colors
	Background() compat.AdaptiveColor        // Radix 1
	BackgroundPanel() compat.AdaptiveColor   // Radix 2
	BackgroundElement() compat.AdaptiveColor // Radix 3

	// Border colors
	BorderSubtle() compat.AdaptiveColor // Radix 6
	Border() compat.AdaptiveColor       // Radix 7
	BorderActive() compat.AdaptiveColor // Radix 8

	// Brand colors
	Primary() compat.AdaptiveColor // Radix 9
	Secondary() compat.AdaptiveColor
	Accent() compat.AdaptiveColor

	// Text colors
	TextMuted() compat.AdaptiveColor // Radix 11
	Text() compat.AdaptiveColor      // Radix 12

	// Status colors
	Error() compat.AdaptiveColor
	Warning() compat.AdaptiveColor
	Success() compat.AdaptiveColor
	Info() compat.AdaptiveColor
}

// BaseTheme provides a default implementation of the Theme interface
// that can be embedded in concrete theme implementations.
type BaseTheme struct {
	// Background colors
	BackgroundColor        compat.AdaptiveColor
	BackgroundPanelColor   compat.AdaptiveColor
	BackgroundElementColor compat.AdaptiveColor

	// Border colors
	BorderSubtleColor compat.AdaptiveColor
	BorderColor       compat.AdaptiveColor
	BorderActiveColor compat.AdaptiveColor

	// Brand colors
	PrimaryColor   compat.AdaptiveColor
	SecondaryColor compat.AdaptiveColor
	AccentColor    compat.AdaptiveColor

	// Text colors
	TextMutedColor compat.AdaptiveColor
	TextColor      compat.AdaptiveColor

	// Status colors
	ErrorColor   compat.AdaptiveColor
	WarningColor compat.AdaptiveColor
	SuccessColor compat.AdaptiveColor
	InfoColor    compat.AdaptiveColor
}

// Implement the Theme interface for BaseTheme
func (t *BaseTheme) Primary() compat.AdaptiveColor   { return t.PrimaryColor }
func (t *BaseTheme) Secondary() compat.AdaptiveColor { return t.SecondaryColor }
func (t *BaseTheme) Accent() compat.AdaptiveColor    { return t.AccentColor }

func (t *BaseTheme) Error() compat.AdaptiveColor   { return t.ErrorColor }
func (t *BaseTheme) Warning() compat.AdaptiveColor { return t.WarningColor }
func (t *BaseTheme) Success() compat.AdaptiveColor { return t.SuccessColor }
func (t *BaseTheme) Info() compat.AdaptiveColor    { return t.InfoColor }

func (t *BaseTheme) Text() compat.AdaptiveColor      { return t.TextColor }
func (t *BaseTheme) TextMuted() compat.AdaptiveColor { return t.TextMutedColor }

func (t *BaseTheme) Background() compat.AdaptiveColor        { return t.BackgroundColor }
func (t *BaseTheme) BackgroundPanel() compat.AdaptiveColor   { return t.BackgroundPanelColor }
func (t *BaseTheme) BackgroundElement() compat.AdaptiveColor { return t.BackgroundElementColor }

func (t *BaseTheme) Border() compat.AdaptiveColor       { return t.BorderColor }
func (t *BaseTheme) BorderActive() compat.AdaptiveColor { return t.BorderActiveColor }
func (t *BaseTheme) BorderSubtle() compat.AdaptiveColor { return t.BorderSubtleColor }
