package timer

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	timer := New()

	assert.Equal(t, Work, timer.SessionType, "initial session should be Work")
	assert.Equal(t, WorkDuration, timer.Duration, "initial duration should be WorkDuration")
	assert.Equal(t, WorkDuration, timer.Remaining, "initial remaining should be WorkDuration")
	assert.False(t, timer.Running, "timer should not be running initially")
	assert.Equal(t, 0, timer.PomodoroCount, "initial pomodoro count should be 0")
	assert.Equal(t, 0, timer.TotalPomodoros, "initial total pomodoros should be 0")
}

func TestStart(t *testing.T) {
	timer := New()
	assert.False(t, timer.Running, "timer should not be running initially")

	timer.Start()
	assert.True(t, timer.Running, "timer should be running after Start()")
}

func TestPause(t *testing.T) {
	timer := New()
	timer.Start()
	assert.True(t, timer.Running, "timer should be running after Start()")

	timer.Pause()
	assert.False(t, timer.Running, "timer should not be running after Pause()")
}

func TestToggle(t *testing.T) {
	timer := New()

	// Toggle from stopped to running
	assert.False(t, timer.Running)
	timer.Toggle()
	assert.True(t, timer.Running, "toggle should start timer")

	// Toggle from running to stopped
	timer.Toggle()
	assert.False(t, timer.Running, "toggle should stop timer")

	// Toggle again
	timer.Toggle()
	assert.True(t, timer.Running, "toggle should start timer again")
}

func TestReset(t *testing.T) {
	timer := New()
	timer.Start()
	timer.Remaining = 10 * time.Minute // Simulate time passing

	timer.Reset()

	assert.Equal(t, timer.Duration, timer.Remaining, "reset should restore remaining to duration")
	assert.False(t, timer.Running, "reset should pause the timer")
}

func TestTick(t *testing.T) {
	tests := []struct {
		name            string
		running         bool
		remaining       time.Duration
		expectedRemain  time.Duration
		shouldDecrement bool
	}{
		{
			name:            "running timer decrements",
			running:         true,
			remaining:       5 * time.Minute,
			expectedRemain:  5*time.Minute - time.Second,
			shouldDecrement: true,
		},
		{
			name:            "paused timer does not decrement",
			running:         false,
			remaining:       5 * time.Minute,
			expectedRemain:  5 * time.Minute,
			shouldDecrement: false,
		},
		{
			name:            "complete timer does not decrement",
			running:         true,
			remaining:       0,
			expectedRemain:  0,
			shouldDecrement: false,
		},
		{
			name:            "timer at exactly zero does not go negative",
			running:         true,
			remaining:       0,
			expectedRemain:  0,
			shouldDecrement: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timer := New()
			timer.Running = tt.running
			timer.Remaining = tt.remaining

			timer.Tick()

			assert.Equal(t, tt.expectedRemain, timer.Remaining)
		})
	}
}

func TestProgress(t *testing.T) {
	tests := []struct {
		name             string
		duration         time.Duration
		remaining        time.Duration
		expectedProgress float64
	}{
		{
			name:             "0% progress at start",
			duration:         25 * time.Minute,
			remaining:        25 * time.Minute,
			expectedProgress: 0.0,
		},
		{
			name:             "50% progress at halfway",
			duration:         10 * time.Minute,
			remaining:        5 * time.Minute,
			expectedProgress: 0.5,
		},
		{
			name:             "100% progress when complete",
			duration:         25 * time.Minute,
			remaining:        0,
			expectedProgress: 1.0,
		},
		{
			name:             "25% progress",
			duration:         20 * time.Minute,
			remaining:        15 * time.Minute,
			expectedProgress: 0.25,
		},
		{
			name:             "zero duration returns 0",
			duration:         0,
			remaining:        0,
			expectedProgress: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timer := New()
			timer.Duration = tt.duration
			timer.Remaining = tt.remaining

			progress := timer.Progress()

			assert.InDelta(t, tt.expectedProgress, progress, 0.001)
		})
	}
}

func TestIsComplete(t *testing.T) {
	tests := []struct {
		name      string
		remaining time.Duration
		expected  bool
	}{
		{
			name:      "not complete with time remaining",
			remaining: 5 * time.Minute,
			expected:  false,
		},
		{
			name:      "complete at zero",
			remaining: 0,
			expected:  true,
		},
		{
			name:      "complete when negative",
			remaining: -time.Second,
			expected:  true,
		},
		{
			name:      "not complete with one second",
			remaining: time.Second,
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timer := New()
			timer.Remaining = tt.remaining

			result := timer.IsComplete()

			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCompleteSession(t *testing.T) {
	tests := []struct {
		name                   string
		initialSession         SessionType
		initialPomodoroCount   int
		initialTotalPomodoros  int
		expectedSession        SessionType
		expectedDuration       time.Duration
		expectedPomodoroCount  int
		expectedTotalPomodoros int
	}{
		{
			name:                   "work to short break (pomodoro 1)",
			initialSession:         Work,
			initialPomodoroCount:   0,
			initialTotalPomodoros:  0,
			expectedSession:        ShortBreak,
			expectedDuration:       ShortBreakDuration,
			expectedPomodoroCount:  1,
			expectedTotalPomodoros: 1,
		},
		{
			name:                   "work to short break (pomodoro 2)",
			initialSession:         Work,
			initialPomodoroCount:   1,
			initialTotalPomodoros:  1,
			expectedSession:        ShortBreak,
			expectedDuration:       ShortBreakDuration,
			expectedPomodoroCount:  2,
			expectedTotalPomodoros: 2,
		},
		{
			name:                   "work to short break (pomodoro 3)",
			initialSession:         Work,
			initialPomodoroCount:   2,
			initialTotalPomodoros:  2,
			expectedSession:        ShortBreak,
			expectedDuration:       ShortBreakDuration,
			expectedPomodoroCount:  3,
			expectedTotalPomodoros: 3,
		},
		{
			name:                   "work to long break (pomodoro 4)",
			initialSession:         Work,
			initialPomodoroCount:   3,
			initialTotalPomodoros:  3,
			expectedSession:        LongBreak,
			expectedDuration:       LongBreakDuration,
			expectedPomodoroCount:  0,
			expectedTotalPomodoros: 4,
		},

		{
			name:                   "short break to work",
			initialSession:         ShortBreak,
			initialPomodoroCount:   1,
			initialTotalPomodoros:  1,
			expectedSession:        Work,
			expectedDuration:       WorkDuration,
			expectedPomodoroCount:  1,
			expectedTotalPomodoros: 1,
		},
		{
			name:                   "long break to work",
			initialSession:         LongBreak,
			initialPomodoroCount:   0,
			initialTotalPomodoros:  4,
			expectedSession:        Work,
			expectedDuration:       WorkDuration,
			expectedPomodoroCount:  0,
			expectedTotalPomodoros: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timer := New()
			timer.SessionType = tt.initialSession
			timer.Duration = WorkDuration // Set some initial duration
			timer.PomodoroCount = tt.initialPomodoroCount
			timer.TotalPomodoros = tt.initialTotalPomodoros
			timer.Running = true

			timer.CompleteSession()

			assert.Equal(t, tt.expectedSession, timer.SessionType)
			assert.Equal(t, tt.expectedDuration, timer.Duration)
			assert.Equal(t, tt.expectedDuration, timer.Remaining)
			assert.Equal(t, tt.expectedPomodoroCount, timer.PomodoroCount)
			assert.Equal(t, tt.expectedTotalPomodoros, timer.TotalPomodoros)
			assert.False(t, timer.Running, "should be paused after completing session")
		})
	}
}

func TestSkip(t *testing.T) {
	tests := []struct {
		name                   string
		initialSession         SessionType
		initialPomodoroCount   int
		initialTotalPomodoros  int
		expectedSession        SessionType
		expectedPomodoroCount  int
		expectedTotalPomodoros int
	}{
		{
			name:                   "skip work session counts toward cycle",
			initialSession:         Work,
			initialPomodoroCount:   0,
			initialTotalPomodoros:  0,
			expectedSession:        ShortBreak,
			expectedPomodoroCount:  1,
			expectedTotalPomodoros: 1,
		},
		{
			name:                   "skip work session triggers long break",
			initialSession:         Work,
			initialPomodoroCount:   3,
			initialTotalPomodoros:  3,
			expectedSession:        LongBreak,
			expectedPomodoroCount:  0,
			expectedTotalPomodoros: 4,
		},
		{
			name:                   "skip short break",
			initialSession:         ShortBreak,
			initialPomodoroCount:   1,
			initialTotalPomodoros:  1,
			expectedSession:        Work,
			expectedPomodoroCount:  1,
			expectedTotalPomodoros: 1,
		},
		{
			name:                   "skip long break",
			initialSession:         LongBreak,
			initialPomodoroCount:   0,
			initialTotalPomodoros:  4,
			expectedSession:        Work,
			expectedPomodoroCount:  0,
			expectedTotalPomodoros: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timer := New()
			timer.SessionType = tt.initialSession
			timer.PomodoroCount = tt.initialPomodoroCount
			timer.TotalPomodoros = tt.initialTotalPomodoros

			timer.Skip()

			assert.Equal(t, tt.expectedSession, timer.SessionType)
			assert.Equal(t, tt.expectedPomodoroCount, timer.PomodoroCount)
			assert.Equal(t, tt.expectedTotalPomodoros, timer.TotalPomodoros)
			assert.False(t, timer.Running)
		})
	}
}

func TestFormatRemaining(t *testing.T) {
	tests := []struct {
		name      string
		remaining time.Duration
		expected  string
	}{
		{
			name:      "25 minutes",
			remaining: 25 * time.Minute,
			expected:  "25:00",
		},
		{
			name:      "5 minutes 30 seconds",
			remaining: 5*time.Minute + 30*time.Second,
			expected:  "05:30",
		},
		{
			name:      "0 minutes 1 second",
			remaining: time.Second,
			expected:  "00:01",
		},
		{
			name:      "0 minutes 0 seconds",
			remaining: 0,
			expected:  "00:00",
		},
		{
			name:      "1 minute exactly",
			remaining: time.Minute,
			expected:  "01:00",
		},
		{
			name:      "9 minutes 59 seconds",
			remaining: 9*time.Minute + 59*time.Second,
			expected:  "09:59",
		},
		{
			name:      "15 minutes",
			remaining: 15 * time.Minute,
			expected:  "15:00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timer := New()
			timer.Remaining = tt.remaining

			result := timer.FormatRemaining()

			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSessionName(t *testing.T) {
	tests := []struct {
		sessionType SessionType
		expected    string
	}{
		{Work, "WORK SESSION"},
		{ShortBreak, "SHORT BREAK"},
		{LongBreak, "LONG BREAK"},
		{SessionType("unknown"), "SESSION"},
	}

	for _, tt := range tests {
		t.Run(string(tt.sessionType), func(t *testing.T) {
			timer := New()
			timer.SessionType = tt.sessionType

			result := timer.SessionName()

			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMinutesRemaining(t *testing.T) {
	tests := []struct {
		remaining time.Duration
		expected  int
	}{
		{25 * time.Minute, 25},
		{5*time.Minute + 30*time.Second, 5},
		{59 * time.Second, 0},
		{0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.remaining.String(), func(t *testing.T) {
			timer := New()
			timer.Remaining = tt.remaining

			assert.Equal(t, tt.expected, timer.MinutesRemaining())
		})
	}
}

func TestSecondsRemaining(t *testing.T) {
	tests := []struct {
		remaining time.Duration
		expected  int
	}{
		{25 * time.Minute, 0},
		{5*time.Minute + 30*time.Second, 30},
		{59 * time.Second, 59},
		{0, 0},
		{1*time.Minute + 1*time.Second, 1},
	}

	for _, tt := range tests {
		t.Run(tt.remaining.String(), func(t *testing.T) {
			timer := New()
			timer.Remaining = tt.remaining

			assert.Equal(t, tt.expected, timer.SecondsRemaining())
		})
	}
}

func TestFullPomodoroCycle(t *testing.T) {
	// Test a complete cycle of 4 pomodoros with breaks
	timer := New()

	require.Equal(t, Work, timer.SessionType)
	require.Equal(t, 0, timer.PomodoroCount)
	require.Equal(t, 0, timer.TotalPomodoros)

	// Complete first work session
	timer.CompleteSession()
	assert.Equal(t, ShortBreak, timer.SessionType)
	assert.Equal(t, 1, timer.TotalPomodoros)
	assert.Equal(t, 1, timer.PomodoroCount)

	// Complete short break
	timer.CompleteSession()
	assert.Equal(t, Work, timer.SessionType)
	assert.Equal(t, 1, timer.PomodoroCount)

	// Complete second work session
	timer.CompleteSession()
	assert.Equal(t, ShortBreak, timer.SessionType)
	assert.Equal(t, 2, timer.TotalPomodoros)
	assert.Equal(t, 2, timer.PomodoroCount)

	// Complete short break
	timer.CompleteSession()
	assert.Equal(t, Work, timer.SessionType)
	assert.Equal(t, 2, timer.PomodoroCount)

	// Complete third work session
	timer.CompleteSession()
	assert.Equal(t, ShortBreak, timer.SessionType)
	assert.Equal(t, 3, timer.TotalPomodoros)
	assert.Equal(t, 3, timer.PomodoroCount)

	// Complete short break
	timer.CompleteSession()
	assert.Equal(t, Work, timer.SessionType)
	assert.Equal(t, 3, timer.PomodoroCount)

	// Complete fourth work session - should trigger long break
	timer.CompleteSession()
	assert.Equal(t, LongBreak, timer.SessionType)
	assert.Equal(t, LongBreakDuration, timer.Duration)
	assert.Equal(t, 4, timer.TotalPomodoros)
	assert.Equal(t, 0, timer.PomodoroCount)

	// Complete long break - back to work, reset count
	timer.CompleteSession()
	assert.Equal(t, Work, timer.SessionType)
	assert.Equal(t, 0, timer.PomodoroCount)
	assert.Equal(t, 4, timer.TotalPomodoros)
}

func TestNextSession(t *testing.T) {
	tests := []struct {
		name                   string
		initialSession         SessionType
		initialPomodoroCount   int
		expectedSession        SessionType
		expectedPomodoroCount  int
		expectedTotalPomodoros int
	}{
		{
			name:                   "work to short break (not at 4)",
			initialSession:         Work,
			initialPomodoroCount:   0,
			expectedSession:        ShortBreak,
			expectedPomodoroCount:  1,
			expectedTotalPomodoros: 1, // NextSession increments total
		},
		{
			name:                   "work to long break (at 4)",
			initialSession:         Work,
			initialPomodoroCount:   3,
			expectedSession:        LongBreak,
			expectedPomodoroCount:  0,
			expectedTotalPomodoros: 1,
		},
		{
			name:                   "short break to work",
			initialSession:         ShortBreak,
			initialPomodoroCount:   1,
			expectedSession:        Work,
			expectedPomodoroCount:  1,
			expectedTotalPomodoros: 0,
		},
		{
			name:                   "long break to work",
			initialSession:         LongBreak,
			initialPomodoroCount:   0,
			expectedSession:        Work,
			expectedPomodoroCount:  0,
			expectedTotalPomodoros: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timer := New()
			timer.SessionType = tt.initialSession
			timer.PomodoroCount = tt.initialPomodoroCount
			timer.TotalPomodoros = 0
			timer.Running = true

			timer.NextSession()

			assert.Equal(t, tt.expectedSession, timer.SessionType)
			assert.Equal(t, tt.expectedPomodoroCount, timer.PomodoroCount)
			assert.Equal(t, tt.expectedTotalPomodoros, timer.TotalPomodoros)
			assert.False(t, timer.Running)
		})
	}
}
