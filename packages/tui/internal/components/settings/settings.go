// ABOUTME: Settings component for configuring Ritual
// ABOUTME: Manages API keys, MCP servers, and other configuration

package settings

import (
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/jem-computer/ritual/tui/internal/api"
)

type Model struct {
	client *api.Client
}

func New(client *api.Client) Model {
	return Model{
		client: client,
	}
}

func (m Model) Init() (tea.Model, tea.Cmd) {
	return m, nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	return "Settings (TODO)"
}
