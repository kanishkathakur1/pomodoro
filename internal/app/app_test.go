package app

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kanishkathakur1/pomodoro/internal/config"
	"github.com/kanishkathakur1/pomodoro/internal/notify"
	"github.com/kanishkathakur1/pomodoro/internal/timer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	// Disable notifications during tests
	notify.SetNotifyFuncsForTesting(
		func(title, message string, icon any) error { return nil },
		func() {},
	)
}

func newTestModel() Model {
	cfg := config.DefaultConfig()
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

func TestNew(t *testing.T) {
	m := New()

	assert.NotNil(t, m.Timer)
	assert.NotNil(t, m.Config)
	assert.NotNil(t, m.Notifier)
	assert.Equal(t, ViewSplash, m.CurrentView)
	assert.Equal(t, 80, m.Width)
	assert.Equal(t, 24, m.Height)
	assert.False(t, m.ShowHelp)
	assert.Equal(t, 0, m.SplashFrame)
	assert.False(t, m.FlashActive)
}

func TestInit(t *testing.T) {
	m := newTestModel()
	cmd := m.Init()

	// Init should return a batch command
	assert.NotNil(t, cmd)
}

func TestUpdate_WindowSizeMsg(t *testing.T) {
	m := newTestModel()
	msg := tea.WindowSizeMsg{Width: 120, Height: 40}

	result, cmd := m.Update(msg)
	model := result.(Model)

	assert.Equal(t, 120, model.Width)
	assert.Equal(t, 40, model.Height)
	assert.Nil(t, cmd)
}

func TestUpdate_SplashTickMsg(t *testing.T) {
	m := newTestModel()
	m.CurrentView = ViewSplash
	m.SplashFrame = 0

	msg := SplashTickMsg(time.Now())
	result, cmd := m.Update(msg)
	model := result.(Model)

	assert.Equal(t, 1, model.SplashFrame)
	assert.NotNil(t, cmd, "should return another splash tick command")
}

func TestUpdate_SplashTickMsg_TransitionsToTimer(t *testing.T) {
	m := newTestModel()
	m.CurrentView = ViewSplash
	m.SplashFrame = 7 // One before transition

	msg := SplashTickMsg(time.Now())
	result, _ := m.Update(msg)
	model := result.(Model)

	assert.Equal(t, 8, model.SplashFrame)
	assert.Equal(t, ViewTimer, model.CurrentView, "should transition to timer after 8 frames")
}

func TestUpdate_SplashTickMsg_IgnoredWhenNotInSplash(t *testing.T) {
	m := newTestModel()
	m.CurrentView = ViewTimer
	m.SplashFrame = 0

	msg := SplashTickMsg(time.Now())
	result, cmd := m.Update(msg)
	model := result.(Model)

	assert.Equal(t, 0, model.SplashFrame, "should not increment when not in splash view")
	assert.Nil(t, cmd)
}

func TestUpdate_TickMsg(t *testing.T) {
	m := newTestModel()
	m.CurrentView = ViewTimer
	m.Timer.Start()
	initialRemaining := m.Timer.Remaining

	msg := TickMsg(time.Now())
	result, cmd := m.Update(msg)
	model := result.(Model)

	assert.Less(t, model.Timer.Remaining, initialRemaining, "timer should decrement")
	assert.NotNil(t, cmd, "should return another tick command")
}

func TestUpdate_TickMsg_IgnoredWhenPaused(t *testing.T) {
	m := newTestModel()
	m.CurrentView = ViewTimer
	m.Timer.Pause()
	initialRemaining := m.Timer.Remaining

	msg := TickMsg(time.Now())
	result, cmd := m.Update(msg)
	model := result.(Model)

	assert.Equal(t, initialRemaining, model.Timer.Remaining, "timer should not decrement when paused")
	assert.Nil(t, cmd)
}

func TestUpdate_TickMsg_IgnoredWhenNotInTimerView(t *testing.T) {
	m := newTestModel()
	m.CurrentView = ViewSplash
	m.Timer.Start()
	initialRemaining := m.Timer.Remaining

	msg := TickMsg(time.Now())
	result, cmd := m.Update(msg)
	model := result.(Model)

	assert.Equal(t, initialRemaining, model.Timer.Remaining)
	assert.Nil(t, cmd)
}

func TestUpdate_TickMsg_SessionComplete(t *testing.T) {
	m := newTestModel()
	m.CurrentView = ViewTimer
	m.Timer.Start()
	m.Timer.Remaining = time.Second // About to complete

	msg := TickMsg(time.Now())
	result, _ := m.Update(msg)
	model := result.(Model)

	assert.Equal(t, ViewComplete, model.CurrentView, "should transition to complete view")
}

func TestUpdate_FlashEndMsg(t *testing.T) {
	m := newTestModel()
	m.FlashActive = true

	msg := FlashEndMsg{}
	result, cmd := m.Update(msg)
	model := result.(Model)

	assert.False(t, model.FlashActive)
	assert.Nil(t, cmd)
}

func TestHandleKey_SplashScreen(t *testing.T) {
	m := newTestModel()
	m.CurrentView = ViewSplash

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
	result, _ := m.Update(msg)
	model := result.(Model)

	assert.Equal(t, ViewTimer, model.CurrentView, "any key should transition from splash to timer")
}

func TestHandleKey_HelpToggle(t *testing.T) {
	m := newTestModel()
	m.CurrentView = ViewTimer
	assert.False(t, m.ShowHelp)

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}
	result, _ := m.Update(msg)
	model := result.(Model)

	assert.True(t, model.ShowHelp, "? should toggle help on")

	// Toggle off
	result, _ = model.Update(msg)
	model = result.(Model)
	assert.False(t, model.ShowHelp, "? again should toggle help off")
}

func TestHandleKey_HelpDismiss(t *testing.T) {
	m := newTestModel()
	m.CurrentView = ViewTimer
	m.ShowHelp = true

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
	result, _ := m.Update(msg)
	model := result.(Model)

	assert.False(t, model.ShowHelp, "any key should dismiss help")
}

func TestHandleKey_Quit(t *testing.T) {
	m := newTestModel()
	m.CurrentView = ViewTimer

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	_, cmd := m.Update(msg)

	// Should return tea.Quit
	assert.NotNil(t, cmd)
}

func TestHandleKey_Toggle(t *testing.T) {
	m := newTestModel()
	m.CurrentView = ViewTimer
	assert.False(t, m.Timer.Running)

	msg := tea.KeyMsg{Type: tea.KeySpace}
	result, cmd := m.Update(msg)
	model := result.(Model)

	assert.True(t, model.Timer.Running, "space should start timer")
	assert.NotNil(t, cmd, "should return tick command")

	// Toggle off
	result, cmd = model.Update(msg)
	model = result.(Model)
	assert.False(t, model.Timer.Running, "space again should pause timer")
	assert.Nil(t, cmd)
}

func TestHandleKey_Skip(t *testing.T) {
	m := newTestModel()
	m.CurrentView = ViewTimer

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}}
	result, _ := m.Update(msg)
	model := result.(Model)

	assert.Equal(t, ViewComplete, model.CurrentView, "s should skip to complete view")
}

func TestHandleKey_Reset(t *testing.T) {
	m := newTestModel()
	m.CurrentView = ViewTimer
	m.Timer.Start()
	m.Timer.Remaining = 10 * time.Minute

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}}
	result, _ := m.Update(msg)
	model := result.(Model)

	assert.Equal(t, m.Timer.Duration, model.Timer.Remaining, "r should reset timer")
	assert.False(t, model.Timer.Running, "r should pause timer")
}

func TestHandleKey_NotifyToggle(t *testing.T) {
	m := newTestModel()
	m.CurrentView = ViewTimer
	initialFlash := m.Config.Notifications.VisualFlash

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}}
	result, _ := m.Update(msg)
	model := result.(Model)

	assert.NotEqual(t, initialFlash, model.Config.Notifications.VisualFlash, "n should toggle notifications")
}

func TestHandleKey_CompleteView_Toggle(t *testing.T) {
	m := newTestModel()
	m.CurrentView = ViewComplete

	msg := tea.KeyMsg{Type: tea.KeySpace}
	result, _ := m.Update(msg)
	model := result.(Model)

	assert.Equal(t, ViewTimer, model.CurrentView, "space in complete view should go to timer")
}

func TestView_Splash(t *testing.T) {
	m := newTestModel()
	m.CurrentView = ViewSplash

	view := m.View()

	assert.NotEmpty(t, view)
	assert.Contains(t, view, "Focus. Flow. Flourish.")
}

func TestView_Timer(t *testing.T) {
	m := newTestModel()
	m.CurrentView = ViewTimer

	view := m.View()

	assert.NotEmpty(t, view)
	assert.Contains(t, view, "WORK SESSION")
}

func TestView_Complete(t *testing.T) {
	m := newTestModel()
	m.CurrentView = ViewComplete
	m.Timer.SessionType = timer.ShortBreak // After completing work

	view := m.View()

	assert.NotEmpty(t, view)
	assert.Contains(t, view, "complete")
}

func TestView_Help(t *testing.T) {
	m := newTestModel()
	m.CurrentView = ViewTimer
	m.ShowHelp = true

	view := m.View()

	assert.NotEmpty(t, view)
	assert.Contains(t, view, "Keyboard Shortcuts")
}

func TestView_Flash(t *testing.T) {
	m := newTestModel()
	m.FlashActive = true

	view := m.View()

	assert.NotEmpty(t, view)
	assert.Contains(t, view, "SESSION COMPLETE!")
}

func TestGetCompletedSession(t *testing.T) {
	tests := []struct {
		name           string
		currentSession timer.SessionType
		pomodoroCount  int
		expected       timer.SessionType
	}{
		{
			name:           "now in short break, completed work",
			currentSession: timer.ShortBreak,
			pomodoroCount:  1,
			expected:       timer.Work,
		},
		{
			name:           "now in long break, completed work",
			currentSession: timer.LongBreak,
			pomodoroCount:  0,
			expected:       timer.Work,
		},
		{
			name:           "now in work (count 1), completed long break",
			currentSession: timer.Work,
			pomodoroCount:  1,
			expected:       timer.LongBreak,
		},
		{
			name:           "now in work (count 2), completed short break",
			currentSession: timer.Work,
			pomodoroCount:  2,
			expected:       timer.ShortBreak,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmr := timer.New()
			tmr.SessionType = tt.currentSession
			tmr.PomodoroCount = tt.pomodoroCount

			result := getCompletedSession(tmr)

			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestHandleSessionComplete(t *testing.T) {
	m := newTestModel()
	m.CurrentView = ViewTimer
	m.Timer.SessionType = timer.Work
	m.Timer.Remaining = 0

	result, cmd := m.handleSessionComplete()
	model := result.(Model)

	assert.Equal(t, ViewComplete, model.CurrentView)
	assert.Equal(t, timer.ShortBreak, model.Timer.SessionType, "should transition to next session")

	// Flash command should be returned if visual flash is enabled
	if m.Config.Notifications.VisualFlash {
		assert.True(t, model.FlashActive)
		assert.NotNil(t, cmd)
	}
}

func TestHandleSessionComplete_NoFlashWhenDisabled(t *testing.T) {
	m := newTestModel()
	m.CurrentView = ViewTimer
	m.Timer.SessionType = timer.Work
	m.Config.Notifications.VisualFlash = false

	result, cmd := m.handleSessionComplete()
	model := result.(Model)

	assert.Equal(t, ViewComplete, model.CurrentView)
	assert.False(t, model.FlashActive)
	assert.Nil(t, cmd)
}

func TestViewStateConstants(t *testing.T) {
	// Verify the view state constants are exported and have correct values
	assert.Equal(t, ViewState(0), ViewSplash)
	assert.Equal(t, ViewState(1), ViewTimer)
	assert.Equal(t, ViewState(2), ViewComplete)
}

func TestModel_InitialState(t *testing.T) {
	m := newTestModel()

	require.NotNil(t, m.Timer)
	assert.Equal(t, timer.Work, m.Timer.SessionType)
	assert.Equal(t, 1, m.Timer.PomodoroCount)
	assert.False(t, m.Timer.Running)
}
