# Pomodoro TUI

A beautiful neon/cyberpunk-styled terminal Pomodoro timer written in Go.

## Features

- Large ASCII countdown display
- Animated splash screen
- Progress bar visualization
- Configurable notifications (visual flash, terminal bell, system notifications)
- Responsive terminal sizing
- Standard Pomodoro timing (25/5/15 minutes)
- Session tracking (4 pomodoros before long break)

## Installation

### From Source

Requires Go 1.24+

```bash
git clone https://github.com/kanishkathakur1/pomodoro.git
cd pomodoro
go build -o pomodoro .
```

### Go Install

```bash
go install github.com/kanishkathakur1/pomodoro@latest
```

Ensure your `$(go env GOPATH)/bin` is in `PATH`, then run `pomodoro`.

### Run

```bash
./pomodoro
```

## Usage

Launch the app and use keyboard controls to manage your focus sessions.

### Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `Space` / `Enter` | Start/pause timer |
| `s` | Skip to next session |
| `r` | Reset current timer |
| `n` | Toggle notifications |
| `?` | Toggle help overlay |
| `q` / `Ctrl+C` | Quit |

### Session Flow

1. **Work Session** (25 minutes) - Focus time
2. **Short Break** (5 minutes) - Quick rest
3. Repeat 4 times
4. **Long Break** (15 minutes) - Extended rest
5. Cycle continues

## Configuration

Configuration is stored at `~/.config/pomodoro/config.toml` and is created automatically on first run.

```toml
[notifications]
visual_flash = true        # Screen flash on session complete
terminal_bell = true       # Terminal bell sound
system_notification = true # Desktop notification
```

## Screenshots

<!-- Add screenshots here -->

## Tech Stack

- [Bubbletea](https://github.com/charmbracelet/bubbletea) - Elm-architecture TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Terminal styling
- [Bubbles](https://github.com/charmbracelet/bubbles) - TUI components
- [Beeep](https://github.com/gen2brain/beeep) - Cross-platform notifications

## Development

```bash
# Run tests
go test ./...

# Run with coverage
go test -coverprofile=coverage.out ./...

# Build
go build -o pomodoro .

# Vet
go vet ./...
```

## License

MIT
