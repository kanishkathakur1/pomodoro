package timer

import (
	"fmt"
	"time"
)

// SessionType represents the type of pomodoro session
type SessionType string

const (
	Work       SessionType = "work"
	ShortBreak SessionType = "short_break"
	LongBreak  SessionType = "long_break"
)

// Standard Pomodoro durations
const (
	WorkDuration       = 25 * time.Minute
	ShortBreakDuration = 5 * time.Minute
	LongBreakDuration  = 15 * time.Minute
	PomodorosBeforeLongBreak = 4
)

// Timer represents the pomodoro timer state
type Timer struct {
	SessionType     SessionType
	Duration        time.Duration
	Remaining       time.Duration
	Running         bool
	PomodoroCount   int // Current pomodoro number in the set (1-4)
	TotalPomodoros  int // Total pomodoros completed
}

// New creates a new timer starting with a work session
func New() *Timer {
	return &Timer{
		SessionType:   Work,
		Duration:      WorkDuration,
		Remaining:     WorkDuration,
		Running:       false,
		PomodoroCount: 1,
		TotalPomodoros: 0,
	}
}

// Start begins the timer
func (t *Timer) Start() {
	t.Running = true
}

// Pause pauses the timer
func (t *Timer) Pause() {
	t.Running = false
}

// Toggle switches between running and paused
func (t *Timer) Toggle() {
	t.Running = !t.Running
}

// Reset resets the current session timer
func (t *Timer) Reset() {
	t.Remaining = t.Duration
	t.Running = false
}

// Tick decrements the timer by one second
func (t *Timer) Tick() {
	if t.Running && t.Remaining > 0 {
		t.Remaining -= time.Second
	}
}

// IsComplete returns true if the current session is done
func (t *Timer) IsComplete() bool {
	return t.Remaining <= 0
}

// Progress returns the completion percentage (0.0 to 1.0)
func (t *Timer) Progress() float64 {
	if t.Duration == 0 {
		return 0
	}
	elapsed := t.Duration - t.Remaining
	return float64(elapsed) / float64(t.Duration)
}

// Skip moves to the next session without completing the current one
func (t *Timer) Skip() {
	t.completeSession(false)
}

// CompleteSession handles the transition to the next session
func (t *Timer) CompleteSession() {
	t.completeSession(true)
}

func (t *Timer) completeSession(wasCompleted bool) {
	// If completing a work session, increment pomodoro count
	if t.SessionType == Work && wasCompleted {
		t.TotalPomodoros++
	}

	// Determine next session
	switch t.SessionType {
	case Work:
		if wasCompleted && t.PomodoroCount >= PomodorosBeforeLongBreak {
			// Time for a long break
			t.SessionType = LongBreak
			t.Duration = LongBreakDuration
			t.PomodoroCount = 0 // Reset after long break
		} else {
			// Short break
			t.SessionType = ShortBreak
			t.Duration = ShortBreakDuration
		}
	case ShortBreak:
		// Back to work after short break
		t.SessionType = Work
		t.Duration = WorkDuration
		t.PomodoroCount++
	case LongBreak:
		// Back to work after long break, reset count
		t.SessionType = Work
		t.Duration = WorkDuration
		t.PomodoroCount = 1
	}

	t.Remaining = t.Duration
	t.Running = false
}

// NextSession prepares for the next session (called after user confirmation)
func (t *Timer) NextSession() {
	// Determine next session
	switch t.SessionType {
	case Work:
		if t.PomodoroCount >= PomodorosBeforeLongBreak {
			// Time for a long break
			t.SessionType = LongBreak
			t.Duration = LongBreakDuration
		} else {
			// Short break
			t.SessionType = ShortBreak
			t.Duration = ShortBreakDuration
		}
		t.TotalPomodoros++
	case ShortBreak:
		// Back to work
		t.SessionType = Work
		t.Duration = WorkDuration
		t.PomodoroCount++
	case LongBreak:
		// Back to work, reset count
		t.SessionType = Work
		t.Duration = WorkDuration
		t.PomodoroCount = 1
	}

	t.Remaining = t.Duration
	t.Running = false
}

// FormatRemaining returns the remaining time as MM:SS
func (t *Timer) FormatRemaining() string {
	minutes := int(t.Remaining.Minutes())
	seconds := int(t.Remaining.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

// SessionName returns a human-readable name for the current session
func (t *Timer) SessionName() string {
	switch t.SessionType {
	case Work:
		return "WORK SESSION"
	case ShortBreak:
		return "SHORT BREAK"
	case LongBreak:
		return "LONG BREAK"
	default:
		return "SESSION"
	}
}

// MinutesRemaining returns the minutes component of the remaining time
func (t *Timer) MinutesRemaining() int {
	return int(t.Remaining.Minutes())
}

// SecondsRemaining returns the seconds component of the remaining time
func (t *Timer) SecondsRemaining() int {
	return int(t.Remaining.Seconds()) % 60
}
