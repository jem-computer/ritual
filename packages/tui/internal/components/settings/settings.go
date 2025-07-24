// ABOUTME: Settings component for configuring Ritual
// ABOUTME: Manages themes, API keys, and other configuration

package settings

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/v2/key"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/jem-computer/ritual/tui/internal/api"
	"github.com/jem-computer/ritual/tui/internal/styles"
	"github.com/jem-computer/ritual/tui/internal/theme"
)

type section int

const (
	sectionTheme section = iota
	// Future sections:
	// sectionAPIKeys
	// sectionMCPServers
	// sectionNotifications
)

type Model struct {
	client *api.Client
	width  int
	height int
	keys   keyMap

	// Current section
	activeSection section

	// Theme settings
	selectedTheme int
	themes        []string
}

type keyMap struct {
	Up     key.Binding
	Down   key.Binding
	Left   key.Binding
	Right  key.Binding
	Select key.Binding
	Back   key.Binding
}

func defaultKeyMap() keyMap {
	return keyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		Left: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("←/h", "prev section"),
		),
		Right: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp("→/l", "next section"),
		),
		Select: key.NewBinding(
			key.WithKeys("enter", " "),
			key.WithHelp("enter/space", "select"),
		),
		Back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back"),
		),
	}
}

func New(client *api.Client) Model {
	return Model{
		client:        client,
		keys:          defaultKeyMap(),
		activeSection: sectionTheme,
		selectedTheme: 0,
		themes: []string{
			"Dracula",
			"Tokyo Night",
			"Catppuccin Mocha",
			"Nord",
			"Gruvbox Dark",
			"Solarized Dark",
			"One Dark",
			"Material",
		},
	}
}

func (m Model) Init() (tea.Model, tea.Cmd) {
	return m, nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Back):
			// TODO: Navigate back to dashboard
			return m, nil

		case key.Matches(msg, m.keys.Up):
			if m.selectedTheme > 0 {
				m.selectedTheme--
			}

		case key.Matches(msg, m.keys.Down):
			if m.selectedTheme < len(m.themes)-1 {
				m.selectedTheme++
			}

		case key.Matches(msg, m.keys.Select):
			// TODO: Apply theme change
			// For now, just log the selection
			return m, nil

		case key.Matches(msg, m.keys.Left):
			// Navigate between sections when we have more
			if m.activeSection > 0 {
				m.activeSection--
			}

		case key.Matches(msg, m.keys.Right):
			// Navigate between sections when we have more
			// For now, we only have one section
		}
	}

	return m, nil
}

func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return ""
	}

	t := theme.CurrentTheme()
	if t == nil {
		return "No theme loaded"
	}

	var s strings.Builder

	// Header
	headerStyle := styles.NewStyle().
		Foreground(t.Primary()).
		Bold(true)

	s.WriteString(headerStyle.Render("> SETTINGS"))
	s.WriteString("\n\n")

	// Section tabs (for future use)
	s.WriteString(m.renderSectionTabs())
	s.WriteString("\n\n")

	// Content based on active section
	switch m.activeSection {
	case sectionTheme:
		s.WriteString(m.renderThemeSection())
	}

	// Help text
	helpStyle := styles.NewStyle().
		Foreground(t.TextMuted()).
		MarginTop(2)

	helpText := "Use ↑/↓ to navigate • Enter/Space to select • ESC to go back"

	// Calculate remaining space
	contentHeight := strings.Count(s.String(), "\n") + 1
	remainingHeight := m.height - contentHeight - 3

	if remainingHeight > 0 {
		s.WriteString(strings.Repeat("\n", remainingHeight))
	}

	s.WriteString("\n")
	s.WriteString(helpStyle.Render(helpText))

	return s.String()
}

func (m Model) renderSectionTabs() string {
	t := theme.CurrentTheme()

	sections := []string{"Theme", "API Keys", "MCP Servers", "Notifications"}
	var tabs []string

	for i, section := range sections {
		style := styles.NewStyle().
			Padding(0, 2).
			MarginRight(1)

		if i == int(m.activeSection) {
			style = style.
				Background(t.Primary()).
				Foreground(t.Background()).
				Bold(true)
		} else if i > 0 {
			// Future sections are disabled
			style = style.
				Foreground(t.TextMuted()).
				Faint(true)
		} else {
			style = style.
				Background(t.BackgroundElement()).
				Foreground(t.Text())
		}

		tabs = append(tabs, style.Render(section))
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, tabs...)
}

func (m Model) renderThemeSection() string {
	t := theme.CurrentTheme()

	var s strings.Builder

	titleStyle := styles.NewStyle().
		Foreground(t.Text()).
		Bold(true).
		MarginBottom(1)

	s.WriteString(titleStyle.Render("Choose Theme"))
	s.WriteString("\n\n")

	// Render theme options as radio buttons
	for i, themeName := range m.themes {
		var radio string
		if i == m.selectedTheme {
			radio = "◉" // Selected
		} else {
			radio = "◯" // Not selected
		}

		optionStyle := styles.NewStyle().
			Foreground(t.Text())

		if i == m.selectedTheme {
			optionStyle = optionStyle.
				Foreground(t.Primary()).
				Bold(true)
		}

		// Add preview colors for each theme
		preview := m.getThemePreview(themeName)

		s.WriteString(fmt.Sprintf("%s %s %s\n",
			optionStyle.Render(radio),
			optionStyle.Render(themeName),
			preview,
		))
	}

	// Current theme indicator
	currentStyle := styles.NewStyle().
		Foreground(t.TextMuted()).
		MarginTop(2)

	s.WriteString("\n")
	s.WriteString(currentStyle.Render(fmt.Sprintf("Current theme: %s", m.themes[0]))) // TODO: Get actual current theme

	return s.String()
}

func (m Model) getThemePreview(themeName string) string {
	// Create a simple color preview for each theme
	// In the future, this will use actual theme colors

	t := theme.CurrentTheme()

	// For now, just use colored blocks to represent each theme
	switch themeName {
	case "Dracula":
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("#bd93f9")).
			Render("████")
	case "Tokyo Night":
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7aa2f7")).
			Render("████")
	case "Catppuccin Mocha":
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("#cba6f7")).
			Render("████")
	case "Nord":
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("#88c0d0")).
			Render("████")
	case "Gruvbox Dark":
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("#fabd2f")).
			Render("████")
	case "Solarized Dark":
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("#268bd2")).
			Render("████")
	case "One Dark":
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("#61afef")).
			Render("████")
	case "Material":
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("#82aaff")).
			Render("████")
	default:
		return lipgloss.NewStyle().
			Foreground(t.Primary()).
			Render("████")
	}
}
