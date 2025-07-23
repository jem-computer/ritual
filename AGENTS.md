# Ritual Codebase Guide for AI Agents

## Build/Lint/Test Commands
```bash
# Development
npm run dev                    # Run both server and TUI concurrently
npm run dev:server            # Run server only (Bun + Hono)
npm run dev:tui              # Run TUI only (Go + Bubbletea)

# Build
npm run build                 # Build both server and TUI
npm run build:server         # Build server to dist/
npm run build:tui            # Build TUI binary to dist/ritual

# Dependencies
npm run install:deps         # Install all dependencies (Bun + Go)
```

## Code Style Guidelines

### TypeScript (Server)
- **Target**: ES2022, strict mode enabled
- **Module**: ESM only (`"type": "module"`)
- **Runtime**: Bun with Hono framework
- **File headers**: Start with 2-line ABOUTME comments
- **Imports**: Use named imports, group by external/internal
- **Error handling**: Use try-catch with proper error responses
- **Types**: Leverage Zod for runtime validation

### Go (TUI)
- **Version**: Go 1.23
- **Framework**: Bubbletea v2 for terminal UI
- **Package structure**: `internal/` for private packages
- **File headers**: Start with 2-line ABOUTME comments
- **Error handling**: Return errors, don't panic
- **Naming**: Use idiomatic Go naming (camelCase for exports)

### General
- **No tests yet**: Project has no test files currently
- **Monorepo**: Uses npm workspaces with packages/server and packages/tui
- **API**: RESTful endpoints at /api/* with JSON responses