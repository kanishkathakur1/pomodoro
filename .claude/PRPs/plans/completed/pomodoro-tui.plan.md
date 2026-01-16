# Feature: Pomodoro TUI

## Summary

Build a beautiful neon/cyberpunk-styled terminal user interface (TUI) for a Pomodoro timer in Go using bubbletea and lipgloss. The app features large ASCII countdown numbers, progress bars, animated splash screen, configurable notifications, and responsive terminal sizing.

## User Story

As a **developer or knowledge worker**
I want to **use a visually appealing terminal-based Pomodoro timer**
So that **I can manage my focus sessions without leaving my terminal environment**

## Problem Statement

Terminal users lack a visually engaging Pomodoro timer that fits their workflow. Existing CLI timers are often plain text-based without visual appeal. This TUI provides a stylish, distraction-free timer with cyberpunk aesthetics.

## Solution Statement

Build a Go TUI application using the Elm-architecture framework bubbletea with lipgloss styling. The app will display ASCII art countdown numbers, progress bars, session state, and provide keyboard controls for timer management. Configuration stored in TOML at XDG-standard paths.

## Metadata

| Field            | Value                                                     |
| ---------------- | --------------------------------------------------------- |
| Type             | NEW_CAPABILITY                                            |
| Complexity       | MEDIUM                                                    |
| Systems Affected | internal/timer, internal/ui, internal/config, internal/notify, main |
| Dependencies     | bubbletea, lipgloss, bubbles, BurntSushi/toml, gen2brain/beeep |
| Estimated Tasks  | 12                                                        |

---

## UX Design

### Before State
```
╔═══════════════════════════════════════════════════════════════════════════════╗
║                              BEFORE STATE                                      ║
╠═══════════════════════════════════════════════════════════════════════════════╣
║                                                                               ║
║   ┌─────────────┐                                                             ║
║   │   No TUI    │                                                             ║
║   │   Exists    │                                                             ║
║   └─────────────┘                                                             ║
║                                                                               ║
║   USER_FLOW: User has no terminal-based Pomodoro timer                        ║
║   PAIN_POINT: Must switch contexts to use web/GUI timer apps                  ║
║   DATA_FLOW: N/A                                                              ║
║                                                                               ║
╚═══════════════════════════════════════════════════════════════════════════════╝
```

### After State
```
╔═══════════════════════════════════════════════════════════════════════════════╗
║                               AFTER STATE                                      ║
╠═══════════════════════════════════════════════════════════════════════════════╣
║                                                                               ║
║   ┌─────────────┐         ┌─────────────┐         ┌─────────────┐            ║
║   │   Launch    │ ──────► │   Splash    │ ──────► │   Timer     │            ║
║   │  `pomodoro` │         │   Screen    │         │    View     │            ║
║   └─────────────┘         └─────────────┘         └─────────────┘            ║
║                                                           │                   ║
║                                                           ▼                   ║
║   ┌─────────────┐         ┌─────────────┐         ┌─────────────┐            ║
║   │    Next     │ ◄────── │  Complete   │ ◄────── │   Timer     │            ║
║   │   Session   │  enter  │    View     │  done   │   Running   │            ║
║   └─────────────┘         └─────────────┘         └─────────────┘            ║
║                                                                               ║
║   USER_FLOW:                                                                  ║
║   1. Run `pomodoro` command                                                   ║
║   2. See animated splash screen (1.5s)                                        ║
║   3. View timer with ASCII numbers + progress bar                             ║
║   4. Press Space/Enter to start/pause                                         ║
║   5. Timer counts down, session completes                                     ║
║   6. Notification fires (visual/bell/system per config)                       ║
║   7. Press Enter to start next session (work → break → work)                  ║
║   8. After 4 pomodoros, get long break                                        ║
║                                                                               ║
║   VALUE_ADD: Beautiful, distraction-free timer in terminal                    ║
║   DATA_FLOW: Config (TOML) → App State → Timer → UI → Notifications           ║
║                                                                               ║
╚═══════════════════════════════════════════════════════════════════════════════╝
```

### Interaction Changes
| Location | Before | After | User Impact |
|----------|--------|-------|-------------|
| Terminal | No app | `pomodoro` launches TUI | Can run Pomodoro from terminal |
| Timer | N/A | Space/Enter toggles | Intuitive pause/resume |
| Session | N/A | Manual transition | User confirms each session start |
| Config | N/A | `~/.config/pomodoro/config.toml` | Customizable notifications |

---

## Mandatory Reading

**CRITICAL: Implementation agent MUST read these files before starting any task:**

| Priority | File | Lines | Why Read This |
|----------|------|-------|---------------|
| P0 | `internal/timer/timer.go` | all | Timer logic already started - understand state |
| P0 | `internal/ui/styles.go` | all | Cyberpunk colors defined - use these |

**External Documentation:**
| Source | Section | Why Needed |
|--------|---------|------------|
| [Bubbletea Tick/Every](https://pkg.go.dev/github.com/charmbracelet/bubbletea) | Tick command | Timer countdown mechanism |
| [Bubbles Progress](https://github.com/charmbracelet/bubbles) | Progress bar | Built-in progress component |
| [Lipgloss Place](https://github.com/charmbracelet/lipgloss) | Placement | Centering content in terminal |
| [BurntSushi/toml](https://github.com/BurntSushi/toml) | Decode/Encode | Config file parsing |
| [gen2brain/beeep](https://github.com/gen2brain/beeep) | Notify | Cross-platform desktop notifications |

---

## Patterns to Mirror

**BUBBLETEA_MODEL_PATTERN:**
```go
// Standard bubbletea Model interface
type Model interface {
    Init() Cmd           // Initial command (start tick, load config)
    Update(Msg) (Model, Cmd)  // Handle messages, return new state + commands
    View() string        // Render UI as string
}
```

**TICK_PATTERN:**
```go
// SOURCE: bubbletea documentation
// Timer tick pattern for countdown
type TickMsg time.Time

func doTick() tea.Cmd {
    return tea.Tick(time.Second, func(t time.Time) tea.Msg {
        return TickMsg(t)
    })
}
```

**KEY_BINDING_PATTERN:**
```go
// SOURCE: bubbles/key documentation
// Define key bindings with help text
keys := keyMap{
    Toggle: key.NewBinding(
        key.WithKeys(" ", "enter"),
        key.WithHelp("space/enter", "start/pause"),
    ),
    Skip:   key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "skip")),
    Reset:  key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "reset")),
    Quit:   key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
}
```

**WINDOW_SIZE_PATTERN:**
```go
// Handle terminal resize
case tea.WindowSizeMsg:
    m.width = msg.Width
    m.height = msg.Height
```

**COLOR_STYLE_PATTERN:**
```go
// SOURCE: internal/ui/styles.go:7-26
// Use existing color definitions
Cyan       = lipgloss.Color("#00FFFF")
Magenta    = lipgloss.Color("#FF00FF")
HotPink    = lipgloss.Color("#FF1493")
```

---

## Files to Change

| File                             | Action | Justification                            |
| -------------------------------- | ------ | ---------------------------------------- |
| `internal/timer/timer.go`        | UPDATE | Fix missing fmt import, add helper methods |
| `internal/ui/styles.go`          | UPDATE | Add any missing styles for new components |
| `internal/ui/ascii.go`           | CREATE | ASCII number rendering for countdown |
| `internal/ui/progress.go`        | CREATE | Custom progress bar component |
| `internal/ui/views.go`           | CREATE | Splash, timer, complete view renderers |
| `internal/ui/help.go`            | CREATE | Help overlay component |
| `internal/config/config.go`      | CREATE | TOML config loading/saving |
| `internal/notify/notify.go`      | CREATE | Notification handling (visual, bell, system) |
| `internal/app/app.go`            | CREATE | Main bubbletea model orchestrating views |
| `internal/app/keys.go`           | CREATE | Key binding definitions |
| `main.go`                        | CREATE | Entry point, tea.Program setup |
| `go.mod`                         | UPDATE | Add dependencies |

---

## NOT Building (Scope Limits)

Explicit exclusions to prevent scope creep:

- **Task/label tracking** - Pure timer, no task management
- **Statistics/history persistence** - No database, no stats file
- **Multiple color themes** - Single cyberpunk theme only
- **Sound files** - Only terminal bell, no custom audio
- **In-app settings UI** - Edit config file directly
- **Custom duration configuration** - Standard 25/5/15 only

---

## Step-by-Step Tasks

Execute in order. Each task is atomic and independently verifiable.

### Task 1: UPDATE `go.mod` - Add dependencies

- **ACTION**: Add required dependencies to go.mod
- **IMPLEMENT**: Add bubbletea, lipgloss, bubbles, toml, beeep
- **COMMAND**:
```bash
go get github.com/charmbracelet/bubbletea@latest
go get github.com/charmbracelet/lipgloss@latest
go get github.com/charmbracelet/bubbles@latest
go get github.com/BurntSushi/toml@latest
go get github.com/gen2brain/beeep@latest
go mod tidy
```
- **VALIDATE**: `go mod tidy && cat go.mod` - verify all deps listed

### Task 2: UPDATE `internal/timer/timer.go` - Fix and enhance

- **ACTION**: Add missing fmt import, enhance Timer struct
- **IMPLEMENT**:
  - Add `import "fmt"`
  - Fix `completeSession` logic bug (line 119 always false)
  - Add `MinutesRemaining()` and `SecondsRemaining()` helpers
- **VALIDATE**: `go build ./internal/timer/...` - must compile

### Task 3: CREATE `internal/config/config.go`

- **ACTION**: Create configuration loading/saving with XDG paths
- **IMPLEMENT**:
```go
package config

type Config struct {
    Notifications NotificationConfig `toml:"notifications"`
}

type NotificationConfig struct {
    VisualFlash        bool `toml:"visual_flash"`
    TerminalBell       bool `toml:"terminal_bell"`
    SystemNotification bool `toml:"system_notification"`
}

func Load() (*Config, error)        // Load from ~/.config/pomodoro/config.toml
func (c *Config) Save() error       // Save config
func DefaultConfig() *Config        // Return sensible defaults
func configPath() (string, error)   // Get XDG config path
```
- **IMPORTS**: `os`, `path/filepath`, `github.com/BurntSushi/toml`
- **GOTCHA**: Use `os.UserConfigDir()` for cross-platform XDG support
- **VALIDATE**: `go build ./internal/config/...`

### Task 4: CREATE `internal/notify/notify.go`

- **ACTION**: Create notification system
- **IMPLEMENT**:
```go
package notify

type Notifier struct {
    config *config.Config
}

func New(cfg *config.Config) *Notifier
func (n *Notifier) Notify(title, message string) error  // Trigger all enabled notifications
func (n *Notifier) VisualFlash() bool                   // Return if visual flash enabled
func terminalBell()                                      // Print \a
func systemNotify(title, msg string) error              // Use beeep.Notify
```
- **IMPORTS**: `github.com/gen2brain/beeep`, `internal/config`
- **GOTCHA**: beeep.Notify icon param can be empty string on macOS
- **VALIDATE**: `go build ./internal/notify/...`

### Task 5: CREATE `internal/ui/ascii.go`

- **ACTION**: Create ASCII art number renderer
- **IMPLEMENT**:
```go
package ui

// Large ASCII digits for timer display
var digits = map[rune][]string{
    '0': {"█████", "█   █", "█   █", "█   █", "█████"},
    '1': {"  █  ", "  █  ", "  █  ", "  █  ", "  █  "},
    // ... all digits 0-9 and ':'
}

func RenderASCII(text string, style lipgloss.Style) string  // Render text as ASCII art
func RenderTime(minutes, seconds int, style lipgloss.Style) string  // Render MM:SS
```
- **PATTERN**: 5 rows per character, space-separated
- **VALIDATE**: `go build ./internal/ui/...`

### Task 6: CREATE `internal/ui/progress.go`

- **ACTION**: Create custom styled progress bar
- **IMPLEMENT**:
```go
package ui

func RenderProgressBar(percent float64, width int) string
// Use filled/empty chars: ████████░░░░░░░░ 65%
// Apply ProgressBarFilled/ProgressBarEmpty styles from styles.go
```
- **MIRROR**: Use `Magenta` for filled, `DarkGray` for empty
- **VALIDATE**: `go build ./internal/ui/...`

### Task 7: CREATE `internal/ui/views.go`

- **ACTION**: Create view rendering functions
- **IMPLEMENT**:
```go
package ui

// Splash screen with animated title
func RenderSplash(frame int) string

// Main timer view with ASCII time, progress, session info
func RenderTimer(t *timer.Timer, width, height int, paused bool) string

// Session complete view with next session prompt
func RenderComplete(completedSession timer.SessionType, nextSession timer.SessionType) string
```
- **PATTERN**: Use lipgloss.Place for centering
- **IMPORTS**: `internal/timer`, `github.com/charmbracelet/lipgloss`
- **VALIDATE**: `go build ./internal/ui/...`

### Task 8: CREATE `internal/ui/help.go`

- **ACTION**: Create help overlay component
- **IMPLEMENT**:
```go
package ui

func RenderHelp() string
// Render boxed help with key bindings:
// space/enter  start/pause
// s            skip session
// r            reset timer
// n            toggle notifications
// ?            toggle help
// q            quit
```
- **STYLE**: Use HelpOverlayStyle from styles.go
- **VALIDATE**: `go build ./internal/ui/...`

### Task 9: CREATE `internal/app/keys.go`

- **ACTION**: Define key bindings
- **IMPLEMENT**:
```go
package app

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
    Toggle key.Binding
    Skip   key.Binding
    Reset  key.Binding
    Notify key.Binding
    Help   key.Binding
    Quit   key.Binding
}

func DefaultKeyMap() KeyMap
```
- **VALIDATE**: `go build ./internal/app/...`

### Task 10: CREATE `internal/app/app.go`

- **ACTION**: Create main bubbletea Model
- **IMPLEMENT**:
```go
package app

type viewState int
const (
    viewSplash viewState = iota
    viewTimer
    viewComplete
)

type Model struct {
    timer      *timer.Timer
    config     *config.Config
    notifier   *notify.Notifier
    keys       KeyMap
    view       viewState
    width      int
    height     int
    showHelp   bool
    splashFrame int
    flashActive bool
}

func New() Model
func (m Model) Init() tea.Cmd       // Load config, start splash timer
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd)
func (m Model) View() string
```
- **MESSAGES**: TickMsg, SplashTickMsg, FlashMsg
- **FLOW**: splash(1.5s) → timer → complete → timer...
- **GOTCHA**: Handle tea.WindowSizeMsg for responsive layout
- **VALIDATE**: `go build ./internal/app/...`

### Task 11: CREATE `main.go`

- **ACTION**: Create entry point
- **IMPLEMENT**:
```go
package main

import (
    "fmt"
    "os"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/kanishkathakur1/pomodoro/internal/app"
)

func main() {
    p := tea.NewProgram(
        app.New(),
        tea.WithAltScreen(),        // Use alternate screen buffer
        tea.WithMouseCellMotion(),  // Enable mouse (optional)
    )
    if _, err := p.Run(); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
}
```
- **VALIDATE**: `go build -o pomodoro .`

### Task 12: INTEGRATION TEST - Full flow verification

- **ACTION**: Manual testing of complete flow
- **TEST STEPS**:
  1. `go build -o pomodoro .`
  2. `./pomodoro` - verify splash screen appears
  3. Wait 1.5s - verify timer view appears
  4. Press Space - verify timer starts counting
  5. Press Space - verify timer pauses
  6. Press `?` - verify help overlay shows
  7. Press `r` - verify timer resets
  8. Press `s` - verify session skips to break
  9. Press `q` - verify clean exit
  10. Check `~/.config/pomodoro/config.toml` exists
- **VALIDATE**: All interactions work as expected

---

## Testing Strategy

### Unit Tests to Write

| Test File                                | Test Cases                 | Validates      |
| ---------------------------------------- | -------------------------- | -------------- |
| `internal/timer/timer_test.go`           | Tick, Progress, NextSession | Timer logic |
| `internal/config/config_test.go`         | Load, Save, Default        | Config parsing |
| `internal/ui/ascii_test.go`              | RenderTime formats         | ASCII rendering |

### Edge Cases Checklist

- [ ] Empty/missing config file → create default
- [ ] Terminal too small (<40 cols) → show minimal view
- [ ] Timer at 0:00 → trigger completion
- [ ] Session 4/4 completes → trigger long break
- [ ] Long break completes → reset to session 1/4
- [ ] Rapid key presses → no state corruption
- [ ] Notification fails → graceful degradation

---

## Validation Commands

### Level 1: STATIC_ANALYSIS

```bash
go build ./... && go vet ./...
```

**EXPECT**: Exit 0, no errors

### Level 2: UNIT_TESTS

```bash
go test ./internal/timer/... ./internal/config/...
```

**EXPECT**: All tests pass

### Level 3: FULL_BUILD

```bash
go build -o pomodoro . && ./pomodoro --help 2>/dev/null || true
```

**EXPECT**: Binary builds successfully

### Level 4: MANUAL_VALIDATION

1. Run `./pomodoro`
2. Verify splash screen appears with animation
3. Timer view shows ASCII countdown "25:00"
4. Space starts timer, counts down
5. Session complete triggers notification
6. Config file created at `~/.config/pomodoro/config.toml`
7. Help overlay (`?`) shows all keybindings
8. Quit (`q`) exits cleanly

---

## Acceptance Criteria

- [ ] Splash screen displays on launch with animation
- [ ] ASCII countdown numbers render correctly
- [ ] Progress bar shows session progress
- [ ] Session counter shows "X/4 until long break"
- [ ] All keyboard controls work (space, s, r, n, ?, q)
- [ ] Help overlay displays and dismisses
- [ ] Config file created with defaults on first run
- [ ] Notifications fire on session complete (per config)
- [ ] Responsive to terminal resize
- [ ] Clean exit with q or Ctrl+C

---

## Completion Checklist

- [ ] Task 1: Dependencies added to go.mod
- [ ] Task 2: Timer.go fixed and enhanced
- [ ] Task 3: Config loading/saving works
- [ ] Task 4: Notifications implemented
- [ ] Task 5: ASCII number rendering works
- [ ] Task 6: Progress bar renders correctly
- [ ] Task 7: All views render (splash, timer, complete)
- [ ] Task 8: Help overlay displays
- [ ] Task 9: Key bindings defined
- [ ] Task 10: Main app model orchestrates views
- [ ] Task 11: main.go launches program
- [ ] Task 12: Full integration test passes
- [ ] Level 1: Static analysis passes
- [ ] Level 2: Unit tests pass
- [ ] Level 3: Build succeeds
- [ ] Level 4: Manual validation complete

---

## Risks and Mitigations

| Risk               | Likelihood | Impact | Mitigation                              |
| ------------------ | ---------- | ------ | --------------------------------------- |
| beeep fails on some systems | MED | LOW | Graceful degradation, log warning |
| Terminal too narrow | LOW | MED | Detect width, show compact view |
| ASCII chars render wrong | LOW | MED | Use safe box-drawing characters |
| Config path varies by OS | MED | MED | Use os.UserConfigDir() |

---

## Notes

### Existing Code Analysis

Two files already created:
- `internal/timer/timer.go` - Good foundation but has bugs:
  - Missing `import "fmt"` for FormatRemaining
  - Line 119: `if t.SessionType == LongBreak` will always be false (just set to Work)
- `internal/ui/styles.go` - Complete cyberpunk color palette, ready to use

### Library Versions

Based on research, recommend:
- bubbletea: latest (supports Model interface)
- lipgloss: latest (v1.x, not v2 beta)
- BurntSushi/toml: v1.x (simpler than pelletier/go-toml)
- beeep: v0.11.x (cross-platform notifications)

### Architecture Decision

Using custom progress bar instead of bubbles/progress for:
1. Tighter style integration with cyberpunk theme
2. Simpler implementation without extra component state
3. Direct use of existing color definitions
