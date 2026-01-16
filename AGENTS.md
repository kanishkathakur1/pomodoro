# Pomodoro TUI - Project Context

## Overview

A beautiful neon/cyberpunk-styled terminal user interface (TUI) for a Pomodoro timer written in Go. Features large ASCII countdown numbers, progress bars, animated splash screen, configurable notifications, and responsive terminal sizing.

## Tech Stack

| Technology | Purpose |
|------------|---------|
| Go 1.24.3 | Language |
| Bubbletea | Elm-architecture TUI framework |
| Lipgloss | Terminal styling and colors |
| Bubbles | TUI components (key bindings) |
| Beeep | Cross-platform desktop notifications |
| BurntSushi/toml | TOML configuration parsing |

## Architecture

```
pomodoro/
├── main.go                      # Entry point - tea.NewProgram setup
├── go.mod / go.sum              # Go module dependencies
└── internal/
    ├── app/
    │   ├── app.go               # Main bubbletea Model (Init, Update, View)
    │   └── keys.go              # Key binding definitions (KeyMap)
    ├── timer/
    │   └── timer.go             # Timer engine - session state, tick logic
    ├── config/
    │   └── config.go            # TOML config loading/saving (XDG paths)
    ├── ui/
    │   ├── styles.go            # Cyberpunk color palette and styles
    │   ├── ascii.go             # ASCII art number rendering
    │   ├── progress.go          # Progress bar component
    │   ├── views.go             # Splash, timer, complete view renderers
    │   └── help.go              # Help overlay component
    └── notify/
        └── notify.go            # Notification system (visual/bell/system)
```

## Key Patterns

### Bubbletea Model Pattern
The app follows the Elm architecture with three main methods:
- `Init()` - Returns initial command (starts splash animation)
- `Update(msg)` - Handles messages, returns new model + commands
- `View()` - Renders UI as a string

### Message Types (internal/app/app.go)
- `TickMsg` - Timer countdown tick (1 second)
- `SplashTickMsg` - Splash animation frame (200ms)
- `FlashMsg` / `FlashEndMsg` - Visual flash notification

### View States
```go
ViewSplash   // Animated splash screen (1.6s)
ViewTimer    // Main countdown view
ViewComplete // Session complete, prompt for next
```

### Session Flow
1. Work (25 min) → Short Break (5 min) → Work → Short Break → ...
2. After 4 work sessions → Long Break (15 min)
3. Cycle repeats

## Module Responsibilities

| Module | Key Types | Purpose |
|--------|-----------|---------|
| `app` | `Model`, `KeyMap` | Orchestrates views, handles input |
| `timer` | `Timer`, `SessionType` | Timer logic, session transitions |
| `config` | `Config`, `NotificationConfig` | Loads/saves `~/.config/pomodoro/config.toml` |
| `ui` | Various render functions | All visual rendering |
| `notify` | `Notifier` | Triggers notifications based on config |

## Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `Space` / `Enter` | Start/pause timer |
| `s` | Skip to next session |
| `r` | Reset current timer |
| `n` | Toggle all notifications |
| `?` | Toggle help overlay |
| `q` / `Ctrl+C` | Quit |

## Configuration

**Path:** `~/.config/pomodoro/config.toml` (XDG standard via `os.UserConfigDir()`)

```toml
[notifications]
visual_flash = true
terminal_bell = true
system_notification = true
```

## Coding Conventions

1. **Imports:** Standard library first, then external deps, then internal packages
2. **Naming:** Exported types are PascalCase, internal functions are camelCase
3. **Error handling:** Graceful degradation - config errors return defaults
4. **Testing:** Each package has `*_test.go` files with table-driven tests

## Scope Boundaries (NOT Building)

- Task/label tracking - Pure timer only
- Statistics/history persistence - No database
- Multiple color themes - Single cyberpunk theme
- Sound files - Only terminal bell
- In-app settings UI - Edit config file directly
- Custom duration configuration - Standard 25/5/15 only

## Validation Commands

```bash
# Build and vet
go build ./... && go vet ./...

# Run tests
go test ./...

# Run with coverage
go test -coverprofile=coverage.out ./...

# Build binary
go build -o pomodoro .

# Run the app
./pomodoro
```

## Important Implementation Details

1. **Timer durations are constants** (`internal/timer/timer.go:18-22`)
2. **Colors are defined in styles.go** - use existing palette, don't add new colors
3. **Session transitions** happen in `timer.completeSession()` and `timer.NextSession()`
4. **Window resizing** is handled via `tea.WindowSizeMsg` in Update
5. **Notifications** check config flags before firing

## Useful File Locations

| Concern | File | Key Lines |
|---------|------|-----------|
| App entry | `main.go` | All |
| Main loop | `internal/app/app.go` | `Update()` method |
| Timer state | `internal/timer/timer.go` | `Timer` struct |
| Color palette | `internal/ui/styles.go` | Lines 6-26 |
| ASCII digits | `internal/ui/ascii.go` | `digits` map |
| Config path | `internal/config/config.go` | `configPath()` |
