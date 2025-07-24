// ABOUTME: Create component for adding new ritual tasks
// ABOUTME: Form-based interface for task creation

package create

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/v2/key"
	"github.com/charmbracelet/bubbles/v2/textarea"
	"github.com/charmbracelet/bubbles/v2/textinput"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/jem-computer/ritual/tui/internal/api"
	"github.com/jem-computer/ritual/tui/internal/styles"
	"github.com/jem-computer/ritual/tui/internal/theme"
)

type state int

const (
	stateForm state = iota
	stateSubmitting
	stateSuccess
	stateError
)

type field int

const (
	fieldName field = iota
	fieldPrompt
	fieldSchedule
	fieldModel
	fieldOutput
)

type Model struct {
	client *api.Client
	state  state
	err    error
	width  int
	height int
	keys   keyMap

	// Form fields
	nameInput     textinput.Model
	promptInput   textarea.Model
	scheduleIndex int
	modelIndex    int
	outputIndex   int

	// Current focused field
	focusedField field

	// Options
	scheduleOptions []string
	modelOptions    []string
	outputOptions   []string
}

type keyMap struct {
	Up     key.Binding
	Down   key.Binding
	Tab    key.Binding
	Submit key.Binding
	Back   key.Binding
}

func defaultKeyMap() keyMap {
	return keyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "shift+tab"),
			key.WithHelp("↑/shift+tab", "prev field"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "tab"),
			key.WithHelp("↓/tab", "next field"),
		),
		Tab: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "next field"),
		),
		Submit: key.NewBinding(
			key.WithKeys("ctrl+s"),
			key.WithHelp("ctrl+s", "submit"),
		),
		Back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back"),
		),
	}
}

func New(client *api.Client) Model {
	// Initialize text input
	nameInput := textinput.New()
	nameInput.Placeholder = "Daily Standup Summary"
	nameInput.CharLimit = 100
	nameInput.Focus()

	// Initialize textarea
	promptInput := textarea.New()
	promptInput.Placeholder = "Summarize my calendar events for today and format as a brief standup update..."
	promptInput.CharLimit = 1000
	promptInput.SetWidth(50)
	promptInput.SetHeight(4)

	return Model{
		client:        client,
		state:         stateForm,
		keys:          defaultKeyMap(),
		nameInput:     nameInput,
		promptInput:   promptInput,
		focusedField:  fieldName,
		scheduleIndex: 0,
		modelIndex:    0,
		outputIndex:   0,
		scheduleOptions: []string{
			"Every day at 9:00 AM",
			"Every Monday at 9:00 AM",
			"Every month on the 1st at 9:00 AM",
			"Every hour",
		},
		modelOptions: []string{
			"GPT-4 (Most capable)",
			"GPT-3.5 Turbo (Faster)",
			"Claude 3 Opus",
			"Claude 3 Sonnet",
		},
		outputOptions: []string{
			"Email to me",
			"SMS",
			"Slack",
			"Discord",
			"Webhook",
		},
	}
}

func (m Model) Init() (tea.Model, tea.Cmd) {
	return m, textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch m.state {
		case stateForm:
			switch {
			case key.Matches(msg, m.keys.Back):
				// TODO: Navigate back to dashboard
				return m, nil

			case key.Matches(msg, m.keys.Submit):
				if m.validate() {
					m.state = stateSubmitting
					return m, m.createTask()
				}

			case key.Matches(msg, m.keys.Up):
				m.focusPrevField()

			case key.Matches(msg, m.keys.Down), key.Matches(msg, m.keys.Tab):
				m.focusNextField()

			default:
				// Handle input based on focused field
				switch m.focusedField {
				case fieldName:
					var cmd tea.Cmd
					m.nameInput, cmd = m.nameInput.Update(msg)
					cmds = append(cmds, cmd)

				case fieldPrompt:
					var cmd tea.Cmd
					m.promptInput, cmd = m.promptInput.Update(msg)
					cmds = append(cmds, cmd)

				case fieldSchedule:
					switch msg.String() {
					case "left", "h":
						if m.scheduleIndex > 0 {
							m.scheduleIndex--
						}
					case "right", "l":
						if m.scheduleIndex < len(m.scheduleOptions)-1 {
							m.scheduleIndex++
						}
					}

				case fieldModel:
					switch msg.String() {
					case "left", "h":
						if m.modelIndex > 0 {
							m.modelIndex--
						}
					case "right", "l":
						if m.modelIndex < len(m.modelOptions)-1 {
							m.modelIndex++
						}
					}

				case fieldOutput:
					switch msg.String() {
					case "left", "h":
						if m.outputIndex > 0 {
							m.outputIndex--
						}
					case "right", "l":
						if m.outputIndex < len(m.outputOptions)-1 {
							m.outputIndex++
						}
					}
				}
			}

		case stateSuccess, stateError:
			switch msg.String() {
			case "enter":
				if m.state == stateSuccess {
					// Reset form for new task
					m.resetForm()
					return m, textinput.Blink
				}
			case "esc":
				// TODO: Navigate back to dashboard
				return m, nil
			}
		}

	case taskCreatedMsg:
		m.state = stateSuccess

	case errorMsg:
		m.state = stateError
		m.err = msg.err
	}

	// Update text inputs
	if m.state == stateForm {
		if m.focusedField == fieldName {
			var cmd tea.Cmd
			m.nameInput, cmd = m.nameInput.Update(msg)
			cmds = append(cmds, cmd)
		} else if m.focusedField == fieldPrompt {
			var cmd tea.Cmd
			m.promptInput, cmd = m.promptInput.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

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

	s.WriteString(headerStyle.Render("> CREATE NEW TASK"))
	s.WriteString("\n\n")

	// Content based on state
	switch m.state {
	case stateForm:
		s.WriteString(m.renderForm())

	case stateSubmitting:
		loadingStyle := styles.NewStyle().
			Foreground(t.TextMuted()).
			Width(m.width-4).
			Height(m.height-10).
			Align(lipgloss.Center, lipgloss.Center)
		s.WriteString(loadingStyle.Render("Creating task..."))

	case stateSuccess:
		successStyle := styles.NewStyle().
			Foreground(t.Success()).
			Width(m.width-4).
			Height(m.height-10).
			Align(lipgloss.Center, lipgloss.Center)
		s.WriteString(successStyle.Render("✓ Task created successfully!\n\nPress ENTER to create another or ESC to go back"))

	case stateError:
		errorStyle := styles.NewStyle().
			Foreground(t.Error()).
			Width(m.width-4).
			Height(m.height-10).
			Align(lipgloss.Center, lipgloss.Center)
		s.WriteString(errorStyle.Render(fmt.Sprintf("Error creating task:\n%v\n\nPress ESC to go back", m.err)))
	}

	return s.String()
}

func (m Model) renderForm() string {
	t := theme.CurrentTheme()

	var s strings.Builder

	// Name field
	s.WriteString(m.renderField("Task Name", m.nameInput.View(), m.focusedField == fieldName))
	s.WriteString("\n\n")

	// Prompt field
	s.WriteString(m.renderField("Prompt", m.promptInput.View(), m.focusedField == fieldPrompt))
	s.WriteString("\n\n")

	// Schedule field
	scheduleView := m.renderSelect(m.scheduleOptions[m.scheduleIndex], m.focusedField == fieldSchedule)
	s.WriteString(m.renderField("Schedule", scheduleView, m.focusedField == fieldSchedule))
	s.WriteString("\n\n")

	// Model field
	modelView := m.renderSelect(m.modelOptions[m.modelIndex], m.focusedField == fieldModel)
	s.WriteString(m.renderField("AI Model", modelView, m.focusedField == fieldModel))
	s.WriteString("\n\n")

	// Output field
	outputView := m.renderSelect(m.outputOptions[m.outputIndex], m.focusedField == fieldOutput)
	s.WriteString(m.renderField("Output Destination", outputView, m.focusedField == fieldOutput))

	// Help text
	helpStyle := styles.NewStyle().
		Foreground(t.TextMuted()).
		MarginTop(2)
	s.WriteString("\n\n")
	s.WriteString(helpStyle.Render("Use ↑/↓ or Tab to navigate • ←/→ to change options • Ctrl+S to submit"))

	return s.String()
}

func (m Model) renderField(label, value string, focused bool) string {
	t := theme.CurrentTheme()

	labelStyle := styles.NewStyle().
		Foreground(t.Text()).
		Bold(true)

	if focused {
		labelStyle = labelStyle.Foreground(t.Primary())
	}

	return labelStyle.Render(label) + "\n" + value
}

func (m Model) renderSelect(value string, focused bool) string {
	t := theme.CurrentTheme()

	style := styles.NewStyle().
		Padding(0, 1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.BorderSubtle()).
		Width(50)

	if focused {
		style = style.BorderForeground(t.Primary())
	}

	return style.Render("◀ " + value + " ▶")
}

func (m *Model) focusNextField() {
	m.focusedField = (m.focusedField + 1) % 5
	m.updateFocus()
}

func (m *Model) focusPrevField() {
	m.focusedField = (m.focusedField + 4) % 5
	m.updateFocus()
}

func (m *Model) updateFocus() {
	m.nameInput.Blur()
	m.promptInput.Blur()

	switch m.focusedField {
	case fieldName:
		m.nameInput.Focus()
	case fieldPrompt:
		m.promptInput.Focus()
	}
}

func (m Model) validate() bool {
	return strings.TrimSpace(m.nameInput.Value()) != "" &&
		strings.TrimSpace(m.promptInput.Value()) != ""
}

func (m *Model) resetForm() {
	m.state = stateForm
	m.err = nil
	m.nameInput.SetValue("")
	m.promptInput.SetValue("")
	m.scheduleIndex = 0
	m.modelIndex = 0
	m.outputIndex = 0
	m.focusedField = fieldName
	m.updateFocus()
}

// Commands

type taskCreatedMsg struct{}

type errorMsg struct {
	err error
}

func (m Model) createTask() tea.Cmd {
	return func() tea.Msg {
		// Map indices to actual values
		scheduleMap := []string{"daily", "weekly", "monthly", "hourly"}
		modelMap := []string{"gpt-4", "gpt-3.5-turbo", "claude-3-opus", "claude-3-sonnet"}
		outputMap := []string{"email", "sms", "slack", "discord", "webhook"}

		task := api.Task{
			Name:     m.nameInput.Value(),
			Prompt:   m.promptInput.Value(),
			Schedule: scheduleMap[m.scheduleIndex],
			Model:    modelMap[m.modelIndex],
			Output:   outputMap[m.outputIndex],
			Status:   "ACTIVE",
			NextRun:  calculateNextRun(scheduleMap[m.scheduleIndex]),
		}

		_, err := m.client.CreateTask(task)
		if err != nil {
			return errorMsg{err: err}
		}

		return taskCreatedMsg{}
	}
}

func calculateNextRun(schedule string) time.Time {
	now := time.Now()

	switch schedule {
	case "hourly":
		return now.Add(time.Hour)
	case "daily":
		// Next day at 9 AM
		next := time.Date(now.Year(), now.Month(), now.Day()+1, 9, 0, 0, 0, now.Location())
		return next
	case "weekly":
		// Next Monday at 9 AM
		daysUntilMonday := (8 - int(now.Weekday())) % 7
		if daysUntilMonday == 0 {
			daysUntilMonday = 7
		}
		next := time.Date(now.Year(), now.Month(), now.Day()+daysUntilMonday, 9, 0, 0, 0, now.Location())
		return next
	case "monthly":
		// First day of next month at 9 AM
		next := time.Date(now.Year(), now.Month()+1, 1, 9, 0, 0, 0, now.Location())
		return next
	default:
		return now.Add(24 * time.Hour)
	}
}
