package app

import (
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kanishkathakur1/pomodoro/internal/config"
	"github.com/kanishkathakur1/pomodoro/internal/notify"
	"github.com/kanishkathakur1/pomodoro/internal/timer"
	"github.com/kanishkathakur1/pomodoro/internal/ui"
)

// ViewState represents the current view
type ViewState int

const (
	ViewSplash ViewState = iota
	ViewTimer
	ViewComplete
)

// Message types
type TickMsg time.Time
type SplashTickMsg time.Time
type FlashMsg struct{}
type FlashEndMsg struct{}

// Model is the main bubbletea model
type Model struct {
	Timer       *timer.Timer
	Config      *config.Config
	Notifier    *notify.Notifier
	Keys        KeyMap
	CurrentView ViewState
	Width       int
	Height      int
	ShowHelp    bool
	SplashFrame int
	FlashActive bool
}

// New creates a new Model
func New() Model {
	cfg, _ := config.Load()
	return Model{
		Timer:       timer.New(),
		Config:      cfg,
		Notifier:    notify.New(cfg),
		Keys:        DefaultKeyMap(),
		CurrentView: ViewSplash,
		Width:       80,
		Height:      24,
	}
}

// Init implements tea.Model
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		splashTick(),
		tea.SetWindowTitle("Pomodoro"),
	)
}

// splashTick creates a tick command for splash animation
func splashTick() tea.Cmd {
	return tea.Tick(200*time.Millisecond, func(t time.Time) tea.Msg {
		return SplashTickMsg(t)
	})
}

// timerTick creates a tick command for the timer
func timerTick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

// flashCmd creates a command to end the flash effect
func flashCmd() tea.Cmd {
	return tea.Tick(500*time.Millisecond, func(t time.Time) tea.Msg {
		return FlashEndMsg{}
	})
}

// Update implements tea.Model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKey(msg)

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil

	case SplashTickMsg:
		if m.CurrentView == ViewSplash {
			m.SplashFrame++
			// After 8 frames (1.6s), transition to timer
			if m.SplashFrame >= 8 {
				m.CurrentView = ViewTimer
				return m, nil
			}
			return m, splashTick()
		}
		return m, nil

	case TickMsg:
		if m.CurrentView == ViewTimer && m.Timer.Running {
			m.Timer.Tick()
			if m.Timer.IsComplete() {
				return m.handleSessionComplete()
			}
			return m, timerTick()
		}
		return m, nil

	case FlashEndMsg:
		m.FlashActive = false
		return m, nil
	}

	return m, nil
}

// handleKey processes keyboard input
func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Handle splash screen - any key transitions to timer
	if m.CurrentView == ViewSplash {
		m.CurrentView = ViewTimer
		return m, nil
	}

	// Handle help toggle in any view
	if key.Matches(msg, m.Keys.Help) {
		m.ShowHelp = !m.ShowHelp
		return m, nil
	}

	// If help is showing, any key hides it
	if m.ShowHelp {
		m.ShowHelp = false
		return m, nil
	}

	// Handle quit
	if key.Matches(msg, m.Keys.Quit) {
		// Save config before quitting
		_ = m.Config.Save()
		return m, tea.Quit
	}

	// Handle view-specific keys
	switch m.CurrentView {
	case ViewTimer:
		return m.handleTimerKey(msg)
	case ViewComplete:
		return m.handleCompleteKey(msg)
	}

	return m, nil
}

// handleTimerKey handles keys in timer view
func (m Model) handleTimerKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.Keys.Toggle):
		m.Timer.Toggle()
		if m.Timer.Running {
			return m, timerTick()
		}
		return m, nil

	case key.Matches(msg, m.Keys.Skip):
		m.Timer.Skip()
		m.CurrentView = ViewComplete
		return m, nil

	case key.Matches(msg, m.Keys.Reset):
		m.Timer.Reset()
		return m, nil

	case key.Matches(msg, m.Keys.Notify):
		m.Notifier.ToggleSystemNotification()
		m.Notifier.ToggleTerminalBell()
		m.Notifier.ToggleVisualFlash()
		return m, nil
	}

	return m, nil
}

// handleCompleteKey handles keys in complete view
func (m Model) handleCompleteKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if key.Matches(msg, m.Keys.Toggle) {
		// Start next session
		m.CurrentView = ViewTimer
		return m, nil
	}
	return m, nil
}

// handleSessionComplete handles timer completion
func (m Model) handleSessionComplete() (tea.Model, tea.Cmd) {
	// Store completed session type before transition
	completedSession := m.Timer.SessionType

	// Send notification
	var title, message string
	switch completedSession {
	case timer.Work:
		title = "Work Session Complete!"
		message = "Time for a break."
	case timer.ShortBreak:
		title = "Break Over!"
		message = "Ready to focus again?"
	case timer.LongBreak:
		title = "Long Break Complete!"
		message = "Great work! Ready for more?"
	}
	_ = m.Notifier.Notify(title, message)

	// Transition to next session
	m.Timer.CompleteSession()
	m.CurrentView = ViewComplete

	// Trigger flash if enabled
	var cmd tea.Cmd
	if m.Notifier.VisualFlash() {
		m.FlashActive = true
		cmd = flashCmd()
	}

	return m, cmd
}

// View implements tea.Model
func (m Model) View() string {
	// Show flash overlay if active
	if m.FlashActive {
		return ui.RenderFlash(m.Width, m.Height)
	}

	// Show help overlay if active
	if m.ShowHelp {
		return ui.RenderHelpCentered(m.Width, m.Height)
	}

	switch m.CurrentView {
	case ViewSplash:
		return ui.RenderSplash(m.SplashFrame, m.Width, m.Height)

	case ViewTimer:
		return ui.RenderTimer(m.Timer, m.Width, m.Height, !m.Timer.Running)

	case ViewComplete:
		// Render complete view centered
		content := ui.RenderComplete(getCompletedSession(m.Timer), m.Timer.SessionType)
		return lipgloss.Place(
			m.Width, m.Height,
			lipgloss.Center, lipgloss.Center,
			content,
		)
	}

	return ""
}

// getCompletedSession infers what session was just completed based on current state
func getCompletedSession(t *timer.Timer) timer.SessionType {
	// After CompleteSession is called, SessionType is updated to the NEXT session
	// So we need to infer what was completed
	switch t.SessionType {
	case timer.Work:
		// If we're now in Work, we just completed a break
		if t.PomodoroCount == 1 {
			return timer.LongBreak
		}
		return timer.ShortBreak
	case timer.ShortBreak, timer.LongBreak:
		// If we're now in a break, we just completed Work
		return timer.Work
	}
	return timer.Work
}
