// ABOUTME: Dashboard component showing all scheduled ritual tasks
// ABOUTME: Displays task list with status, schedule, and quick actions

package dashboard

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/v2/key"
	"github.com/charmbracelet/bubbles/v2/list"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/jem-computer/ritual/tui/internal/api"
	"github.com/jem-computer/ritual/tui/internal/styles"
	"github.com/jem-computer/ritual/tui/internal/theme"
)

// taskItem implements list.Item interface for api.Task
type taskItem struct {
	task api.Task
}

func (i taskItem) FilterValue() string {
	return i.task.Name
}

func (i taskItem) Title() string {
	return i.task.Name
}

func (i taskItem) Description() string {
	status := i.task.Status
	if status == "PAUSED" {
		status = "⏸ PAUSED"
	} else {
		status = "▶ ACTIVE"
	}

	nextRun := "Next: " + i.task.NextRun.Format("Jan 2, 3:04 PM")
	if i.task.LastRun.IsZero() {
		return fmt.Sprintf("%s • %s • Never run", status, nextRun)
	}
	return fmt.Sprintf("%s • %s • Last: %s", status, nextRun, i.task.LastRun.Format("Jan 2, 3:04 PM"))
}

// Custom item delegate for better control over rendering
type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 3 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(taskItem)
	if !ok {
		return
	}

	t := theme.CurrentTheme()
	if t == nil {
		return
	}

	var s strings.Builder

	if index == m.Index() {
		// Selected item - add a bullet point and highlight
		title := styles.NewStyle().
			Foreground(t.Primary()).
			Bold(true).
			Render("✦ " + i.Title())
		desc := styles.NewStyle().
			Foreground(t.TextMuted()).
			PaddingLeft(2).
			MarginBottom(1).
			Render(i.Description())
		s.WriteString(title)
		s.WriteString("\n")
		s.WriteString(desc)
	} else {
		// Normal item
		title := styles.NewStyle().
			Foreground(t.Text()).
			PaddingLeft(2).
			Render(i.Title())
		desc := styles.NewStyle().
			Foreground(t.TextMuted()).
			PaddingLeft(2).
			MarginBottom(1).
			Render(i.Description())
		s.WriteString(title)
		s.WriteString("\n")
		s.WriteString(desc)
	}

	fmt.Fprint(w, s.String())
}

type Model struct {
	client *api.Client
	list   list.Model
	tasks  []api.Task
	width  int
	height int
	keys   keyMap
	err    error
}

type keyMap struct {
	Up     key.Binding
	Down   key.Binding
	Enter  key.Binding
	Delete key.Binding
	Pause  key.Binding
	Retry  key.Binding
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
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "edit"),
		),
		Delete: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "delete"),
		),
		Pause: key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp("p", "pause/resume"),
		),
		Retry: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "retry"),
		),
	}
}

func New(client *api.Client) Model {
	// Use our custom delegate
	delegate := itemDelegate{}

	l := list.New([]list.Item{}, delegate, 0, 0)
	l.Title = ""
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(false) // Disable filtering for now
	l.SetShowHelp(true)
	l.DisableQuitKeybindings()

	// Update list keybindings to match our custom ones
	keys := defaultKeyMap()
	l.KeyMap.CursorUp = keys.Up
	l.KeyMap.CursorDown = keys.Down

	// Add our custom keys to the list's help
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			keys.Enter,
			keys.Pause,
			keys.Delete,
		}
	}
	l.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			keys.Enter,
			keys.Pause,
			keys.Delete,
			keys.Retry,
		}
	}

	return Model{
		client: client,
		list:   l,
		keys:   keys,
	}
}
func (m Model) Init() (tea.Model, tea.Cmd) {
	return m, m.loadTasks
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height - 4 // Account for tab bar and status bar
		// List height: total height - header (2 lines) - NEW TASK button (1 line) - spacing
		listHeight := m.height - 5
		if listHeight < 0 {
			listHeight = 0
		}
		m.list.SetSize(m.width-4, listHeight)

	case tea.KeyMsg:
		// Handle custom keybindings first
		switch {
		case key.Matches(msg, m.keys.Delete):
			if selectedItem, ok := m.list.SelectedItem().(taskItem); ok {
				return m, m.deleteTask(selectedItem.task.ID)
			}

		case key.Matches(msg, m.keys.Pause):
			if selectedItem, ok := m.list.SelectedItem().(taskItem); ok {
				return m, m.toggleTaskStatus(selectedItem.task.ID)
			}

		case key.Matches(msg, m.keys.Retry):
			if m.err != nil {
				m.err = nil
				return m, m.loadTasks
			}
		}

	case tasksLoadedMsg:
		m.tasks = msg.tasks
		// Convert tasks to list items
		items := make([]list.Item, len(m.tasks))
		for i, task := range m.tasks {
			items[i] = taskItem{task: task}
		}
		cmd := m.list.SetItems(items)
		cmds = append(cmds, cmd)

	case taskDeletedMsg:
		return m, m.loadTasks

	case taskUpdatedMsg:
		return m, m.loadTasks

	case errorMsg:
		m.err = msg.err
		// Clear tasks on error
		m.tasks = []api.Task{}
		cmd := m.list.SetItems([]list.Item{})
		cmds = append(cmds, cmd)
	}

	// Update the list
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
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

	s.WriteString(headerStyle.Render("> SCHEDULED TASKS"))
	s.WriteString("\n\n")

	// NEW TASK button at bottom right
	buttonStyle := styles.NewStyle().
		Background(t.Primary()).
		Foreground(t.Background()).
		Padding(0, 2).
		Bold(true)

	s.WriteString(lipgloss.PlaceHorizontal(m.width-4, lipgloss.Left, buttonStyle.Render("+ NEW TASK")))

	if m.err != nil {
		// Error state
		errorStyle := styles.NewStyle().
			Foreground(t.Error()).
			Width(m.width-4).
			Height(m.height-10).
			Align(lipgloss.Center, lipgloss.Center)
		s.WriteString(errorStyle.Render(fmt.Sprintf("Error loading tasks:\n%v\n\nPress [R] to retry", m.err)))
	} else if len(m.tasks) == 0 {
		// Empty state
		emptyStyle := styles.NewStyle().
			Foreground(t.TextMuted()).
			Width(m.width-4).
			Height(m.height-10).
			Align(lipgloss.Center, lipgloss.Center)
		s.WriteString(emptyStyle.Render("No tasks scheduled\n\nPress [C] to create a new task"))
	} else {
		// Render the list
		s.WriteString(m.list.View())
	}

	// Calculate remaining space
	contentHeight := lipgloss.Height(s.String())
	remainingHeight := m.height - contentHeight - 3 // Leave space for NEW TASK button

	if remainingHeight > 0 {
		s.WriteString(strings.Repeat("\n", remainingHeight))
	}

	return s.String()
}

// Commands

type tasksLoadedMsg struct {
	tasks []api.Task
}

type taskDeletedMsg struct{}

type taskUpdatedMsg struct{}

type errorMsg struct {
	err error
}

func (m Model) loadTasks() tea.Msg {
	tasks, err := m.client.GetTasks()
	if err != nil {
		return errorMsg{err: err}
	}

	return tasksLoadedMsg{tasks: tasks}
}

func (m Model) deleteTask(id string) tea.Cmd {
	return func() tea.Msg {
		err := m.client.DeleteTask(id)
		if err != nil {
			return errorMsg{err: err}
		}
		return taskDeletedMsg{}
	}
}

func (m Model) toggleTaskStatus(id string) tea.Cmd {
	return func() tea.Msg {
		// Find the task
		var task *api.Task
		for _, t := range m.tasks {
			if t.ID == id {
				task = &t
				break
			}
		}

		if task == nil {
			return nil
		}

		// Toggle status
		if task.Status == "ACTIVE" {
			task.Status = "PAUSED"
		} else {
			task.Status = "ACTIVE"
		}

		_, err := m.client.UpdateTask(id, *task)
		if err != nil {
			return errorMsg{err: err}
		}
		return taskUpdatedMsg{}
	}
}
