// ABOUTME: Main TUI model and logic for the Ritual terminal interface
// ABOUTME: Implements the Bubbletea Model interface with tab navigation

package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/v2/key"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/jem-computer/ritual/tui/internal/api"
	"github.com/jem-computer/ritual/tui/internal/components/create"
	"github.com/jem-computer/ritual/tui/internal/components/dashboard"
	"github.com/jem-computer/ritual/tui/internal/components/logs"
	"github.com/jem-computer/ritual/tui/internal/components/settings"
	"github.com/jem-computer/ritual/tui/internal/styles"
	"github.com/jem-computer/ritual/tui/internal/theme"
)

type Tab int

const (
	DashboardTab Tab = iota
	CreateTab
	LogsTab
	SettingsTab
)

type Model struct {
	width, height int
	activeTab     Tab
	client        *api.Client
	version       string

	// Components
	dashboard dashboard.Model
	create    create.Model
	logs      logs.Model
	settings  settings.Model

	// Key bindings
	keys keyMap
}

type keyMap struct {
	Tab      key.Binding
	ShiftTab key.Binding
	Quit     key.Binding
	Help     key.Binding
}

func defaultKeyMap() keyMap {
	return keyMap{
		Tab: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "next tab"),
		),
		ShiftTab: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("shift+tab", "previous tab"),
		),
		Quit: key.NewBinding(
			key.WithKeys("ctrl+c", "esc"),
			key.WithHelp("ctrl+c/esc", "quit"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
		),
	}
}

func New(client *api.Client, version string) Model {
	return Model{
		activeTab: DashboardTab,
		client:    client,
		version:   version,
		dashboard: dashboard.New(client),
		create:    create.New(client),
		logs:      logs.New(client),
		settings:  settings.New(client),
		keys:      defaultKeyMap(),
	}
}

func (m Model) Init() (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	dashboardModel, cmd := m.dashboard.Init()
	m.dashboard = dashboardModel.(dashboard.Model)
	cmds = append(cmds, cmd)

	createModel, cmd := m.create.Init()
	m.create = createModel.(create.Model)
	cmds = append(cmds, cmd)

	logsModel, cmd := m.logs.Init()
	m.logs = logsModel.(logs.Model)
	cmds = append(cmds, cmd)

	settingsModel, cmd := m.settings.Init()
	m.settings = settingsModel.(settings.Model)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Update all components with new size
		dashboardModel, cmd := m.dashboard.Update(msg)
		m.dashboard = dashboardModel.(dashboard.Model)
		cmds = append(cmds, cmd)

		createModel, cmd := m.create.Update(msg)
		m.create = createModel.(create.Model)
		cmds = append(cmds, cmd)

		logsModel, cmd := m.logs.Update(msg)
		m.logs = logsModel.(logs.Model)
		cmds = append(cmds, cmd)

		settingsModel, cmd := m.settings.Update(msg)
		m.settings = settingsModel.(settings.Model)
		cmds = append(cmds, cmd)

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, m.keys.Tab):
			m.activeTab = (m.activeTab + 1) % 4

		case key.Matches(msg, m.keys.ShiftTab):
			m.activeTab = (m.activeTab + 3) % 4

		case msg.String() == "d":
			m.activeTab = DashboardTab

		case msg.String() == "c":
			m.activeTab = CreateTab

		case msg.String() == "l":
			m.activeTab = LogsTab

		case msg.String() == "s":
			m.activeTab = SettingsTab
		}
	}

	// Update active component
	switch m.activeTab {
	case DashboardTab:
		newModel, cmd := m.dashboard.Update(msg)
		m.dashboard = newModel.(dashboard.Model)
		cmds = append(cmds, cmd)

	case CreateTab:
		newModel, cmd := m.create.Update(msg)
		m.create = newModel.(create.Model)
		cmds = append(cmds, cmd)

	case LogsTab:
		newModel, cmd := m.logs.Update(msg)
		m.logs = newModel.(logs.Model)
		cmds = append(cmds, cmd)

	case SettingsTab:
		newModel, cmd := m.settings.Update(msg)
		m.settings = newModel.(settings.Model)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return ""
	}

	// Tab bar
	tabBar := m.renderTabBar()
	tabBarHeight := lipgloss.Height(tabBar)

	// Status bar
	statusBar := m.renderStatusBar()
	statusBarHeight := lipgloss.Height(statusBar)

	// Content area with calculated height
	contentHeight := m.height - tabBarHeight - statusBarHeight
	var content string
	switch m.activeTab {
	case DashboardTab:
		content = m.dashboard.View()
	case CreateTab:
		content = m.create.View()
	case LogsTab:
		content = m.logs.View()
	case SettingsTab:
		content = m.settings.View()
	}

	// Ensure content fills the available space
	contentLines := strings.Split(content, "\n")
	if len(contentLines) < contentHeight {
		// Add empty lines to fill the space
		for i := len(contentLines); i < contentHeight; i++ {
			content += "\n"
		}
	}

	// Combine all sections
	return lipgloss.JoinVertical(
		lipgloss.Top,
		tabBar,
		content,
		statusBar,
	)
}

func (m Model) renderTabBar() string {
	t := theme.CurrentTheme()
	if t == nil {
		return "No theme"
	}

	tabs := []string{"Dashboard", "Create", "Logs", "Settings"}
	var renderedTabs []string

	for i, tab := range tabs {
		style := styles.NewStyle().
			Padding(0, 2).
			MarginRight(1)

		if i == int(m.activeTab) {
			style = style.
				Background(t.Primary()).
				Foreground(t.Background()).
				Bold(true)
		} else {
			style = style.
				Background(t.BackgroundElement()).
				Foreground(t.Text())
		}

		// Add keyboard shortcut hint
		shortcut := fmt.Sprintf("[%s]", strings.ToLower(tab[:1]))
		renderedTabs = append(renderedTabs, style.Render(shortcut+tab))
	}

	tabBar := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	return styles.NewStyle().
		Width(m.width).
		Background(t.BackgroundPanel()).
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		BorderForeground(t.BorderSubtle()).
		Render(tabBar)
}

func (m Model) renderStatusBar() string {
	t := theme.CurrentTheme()
	if t == nil {
		return "No theme"
	}

	left := fmt.Sprintf(" Ritual %s", m.version)
	right := "ESC to exit "

	// Calculate padding
	padding := m.width - len(left) - len(right)
	if padding < 0 {
		padding = 0
	}

	status := left + strings.Repeat(" ", padding) + right

	return styles.NewStyle().
		Width(m.width).
		Background(t.BackgroundPanel()).
		Foreground(t.TextMuted()).
		Render(status)
}
