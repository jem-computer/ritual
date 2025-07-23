// ABOUTME: Create component for adding new ritual tasks
// ABOUTME: Form-based interface for task creation

package create

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
	return "Create Task (TODO)"
}
