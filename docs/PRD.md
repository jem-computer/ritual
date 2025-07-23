# Ritual Product Requirements Document

## Overview

Ritual is a terminal-based tool that allows users to create automated LLM "rituals" - scheduled AI tasks that run at specified times, with results delivered via various output methods.

## Core Value Proposition

- Schedule AI tasks to run automatically (e.g., "daily at 8am, summarize today's Things tasks")
- Beautiful, keyboard-driven TUI interface
- Connect to local and remote MCP servers for context
- Multiple output destinations (SMS, email, Slack, etc.)
- Test tasks with one-off runs before scheduling

## User Interface

### Dashboard View

The main screen showing all scheduled tasks with:

- Task name and status (ACTIVE/PAUSED)
- Prompt preview (truncated)
- Schedule details and next run time
- Last run time
- Output destination
- Model being used
- Quick actions via keyboard shortcuts (pause/play, edit, delete)
- "NEW TASK" button prominently displayed

### Create Task Flow

Simple form-based creation with:

- Task name
- Prompt (with expandable text area)
- Schedule dropdown (with common presets + custom cron)
- Model selection dropdown
- Output destination field
- Create/Cancel buttons

### Logs View

Chronological execution history showing:

- Timestamp
- Task name
- Output destination
- Success/Error status
- Error messages when applicable

### Navigation

- Tab-based navigation: [D]ashboard, [C]reate, [L]ogs, [S]ettings
- ESC to exit
- Keyboard shortcuts displayed at bottom

## Technical Architecture

- **Frontend**: Go + Bubbletea for the TUI
- **Backend**: TypeScript/Bun Hono server handling scheduling and LLM calls
- **Communication**: Same pattern as OpenCode (local server spawned by TUI)

## Key Features

### V1.0 Must-Haves

1. **Task Management**
   - Create tasks with name, prompt, schedule
   - List all tasks in dashboard
   - Pause/resume tasks
   - Delete tasks
   - Edit existing tasks

2. **Scheduling**
   - Cron expression support
   - Common presets (daily, weekly, etc.)
   - Display next run time
   - Timezone awareness

3. **Execution**
   - Run tasks on schedule
   - Manual "run now" for testing
   - Basic retry on failure

4. **MCP Server integration**
   - Add MCP servers in `ritual.json` config file

5. **Output Methods**
   - SMS, Email, Slack etc (via MCP integrations)
   - Terminal notification

6. **Models**
   - Any supported through AI SDK / [Models.dev](https://models.dev/)

7. **Logging**
   - Execution history
   - Success/failure tracking
   - Basic error messages

### V1.1 Nice-to-Haves

- MCP server integration
- Webhook outputs
- Advanced scheduling (e.g., "every last Friday")
- Export/import tasks
- Task templates

### Future Considerations

- Web UI companion
- Task dependencies
- Conditional execution
- Variables/parameters in prompts
- Cost tracking
- Team/sharing features

## Design Principles

1. **Keyboard-first**: Everything accessible without mouse
2. **Information density**: Show relevant info without clutter
3. **Fast feedback**: Immediate visual feedback for all actions
4. **Fail gracefully**: Clear error messages, never crash
5. **Beautiful defaults**: Works great out of the box

## Success Metrics

- Users can create and schedule a task in under 60 seconds
- Tasks execute reliably 99%+ of the time
- TUI remains responsive even with 100+ tasks
- Clear indication when tasks fail and why

## Non-Goals for V1

- Complex workflow orchestration
- Multi-step prompts
- User authentication/multi-user
- Cloud hosting (local only)
- Mobile app
- Cost optimization features

## Open Questions

1. Should we support multiple LLM providers in V1 or just OpenAI?
   - use AI SDK and [Model.dev](https://models.dev/) like OpenCode does
2. How do we handle API key management? Environment variables vs config file?
   - do whatever OpenCode does
3. Should task history be pruned automatically or kept indefinitely?
   - keep indefinitely
4. Do we need task categories/folders for organization?
   - no
5. Should we support task templates in V1?
   - not yet

## Example Use Cases

1. **Daily Standup**: "Summarize my Things tasks and Linear tickets, send to Slack"
2. **Weekly Report**: "Each Sunday, Generate productivity report from calendar and git commits"
3. **Content Ideas**: "Every Monday, generate 5 blog post ideas based on trending topics"
4. **Task Reminders**: "Check for overdue tasks and text me if any exist"
5. **Commit Summary**: "End of day, summarize today's commits in team channel"
