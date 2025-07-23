# Ritual

A terminal-based tool for scheduling automated LLM tasks with beautiful TUI interface.

## Architecture

Ritual follows the same client-server architecture as OpenCode:

- **TUI Client** (Go + Bubbletea): Terminal interface for managing tasks
- **Server** (TypeScript/Bun): HTTP API server handling scheduling and LLM calls
- **Communication**: HTTP API + SSE for real-time updates

## Project Structure

```
ritual/
├── packages/
│   ├── tui/                    # Go TUI client
│   │   ├── cmd/ritual/         # Main entry point
│   │   └── internal/           # Internal packages
│   │       ├── api/            # API client
│   │       ├── tui/            # Main TUI model
│   │       └── components/     # UI components
│   │           ├── dashboard/  # Task list view
│   │           ├── create/     # Task creation form
│   │           ├── logs/       # Execution history
│   │           └── settings/   # Configuration
│   └── server/                 # TypeScript server
│       └── src/
│           └── index.ts        # HTTP API server
├── package.json                # Root package.json with scripts
└── README.md                   # This file
```

## Development

### Prerequisites

- Go 1.23+
- Bun 1.0+
- Node.js 20+ (for npm workspaces)

### Setup

1. Install dependencies:
```bash
npm run install:deps
```

2. Run in development mode:
```bash
npm run dev
```

This will start both the server (on port 8080) and the TUI client.

### Building

```bash
npm run build
```

This creates:
- `dist/ritual` - The TUI binary
- `packages/server/dist/` - The server bundle

## Features (TODO)

- [ ] Task scheduling with cron expressions
- [ ] Multiple LLM provider support
- [ ] MCP server integration
- [ ] Output methods (SMS, Email, Slack)
- [ ] Execution history and logs
- [ ] Task templates

## Design Principles

- **Keyboard-first**: Everything accessible without mouse
- **Beautiful defaults**: Works great out of the box
- **Fast feedback**: Immediate visual feedback for all actions
- **Fail gracefully**: Clear error messages, never crash