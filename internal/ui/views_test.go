package ui

import (
	"testing"

	"github.com/kanishkathakur1/pomodoro/internal/timer"
	"github.com/stretchr/testify/assert"
)

func TestRenderSplash(t *testing.T) {
	tests := []struct {
		name   string
		frame  int
		width  int
		height int
	}{
		{"frame 0", 0, 80, 24},
		{"frame 1", 1, 80, 24},
		{"frame 7", 7, 80, 24},
		{"wide screen", 0, 120, 40},
		{"narrow screen", 0, 60, 20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RenderSplash(tt.frame, tt.width, tt.height)

			// Should contain title ASCII art elements
			assert.NotEmpty(t, result)

			// Should contain subtitle
			assert.Contains(t, result, "Focus. Flow. Flourish.")
			assert.Contains(t, result, "Press any key to start...")
		})
	}
}

func TestRenderSplash_AnimatesColors(t *testing.T) {
	// Different frames should produce different outputs due to color animation
	result0 := RenderSplash(0, 80, 24)
	result1 := RenderSplash(1, 80, 24)

	// Both should be valid
	assert.NotEmpty(t, result0)
	assert.NotEmpty(t, result1)

	// They may differ due to color animation
	// (Can't easily verify this without disabling ANSI codes)
}

func TestRenderSplash_ContainsASCIITitle(t *testing.T) {
	result := RenderSplash(0, 100, 30)

	// Should contain some ASCII art characters
	assert.Contains(t, result, "█")
	assert.Contains(t, result, "╗")
}

func TestRenderTimer_WorkSession(t *testing.T) {
	tmr := timer.New()

	result := RenderTimer(tmr, 80, 24, false)

	assert.NotEmpty(t, result)
	assert.Contains(t, result, "WORK SESSION")
	assert.Contains(t, result, "RUNNING")
	assert.Contains(t, result, "Press ? for help")
}

func TestRenderTimer_Paused(t *testing.T) {
	tmr := timer.New()

	result := RenderTimer(tmr, 80, 24, true)

	assert.NotEmpty(t, result)
	assert.Contains(t, result, "PAUSED")
}

func TestRenderTimer_ShortBreak(t *testing.T) {
	tmr := timer.New()
	tmr.SessionType = timer.ShortBreak

	result := RenderTimer(tmr, 80, 24, false)

	assert.NotEmpty(t, result)
	assert.Contains(t, result, "SHORT BREAK")
}

func TestRenderTimer_LongBreak(t *testing.T) {
	tmr := timer.New()
	tmr.SessionType = timer.LongBreak

	result := RenderTimer(tmr, 80, 24, false)

	assert.NotEmpty(t, result)
	assert.Contains(t, result, "LONG BREAK")
}

func TestRenderTimer_ShowsSessionCounter(t *testing.T) {
	tmr := timer.New()
	tmr.PomodoroCount = 2

	result := RenderTimer(tmr, 80, 24, false)

	assert.Contains(t, result, "Pomodoro 2/4")
}

func TestRenderTimer_ShowsNextBreakInfo(t *testing.T) {
	tmr := timer.New()
	tmr.PomodoroCount = 3

	result := RenderTimer(tmr, 80, 24, false)

	assert.Contains(t, result, "Short break next")
}

func TestRenderTimer_ShowsLongBreakNext(t *testing.T) {
	tmr := timer.New()
	tmr.PomodoroCount = 4

	result := RenderTimer(tmr, 80, 24, false)

	assert.Contains(t, result, "Long break next!")
}

func TestRenderTimer_BreakShowsWorkNext(t *testing.T) {
	tmr := timer.New()
	tmr.SessionType = timer.ShortBreak

	result := RenderTimer(tmr, 80, 24, false)

	assert.Contains(t, result, "Work session next")
}

func TestRenderTimer_LongBreakShowsFullCount(t *testing.T) {
	tmr := timer.New()
	tmr.SessionType = timer.LongBreak
	tmr.PomodoroCount = 0

	result := RenderTimer(tmr, 80, 24, false)

	assert.Contains(t, result, "Pomodoro 4/4")
}

func TestRenderTimer_NarrowWidth(t *testing.T) {
	tmr := timer.New()

	result := RenderTimer(tmr, 50, 24, false)

	assert.NotEmpty(t, result)
}

func TestRenderComplete_WorkComplete(t *testing.T) {
	result := RenderComplete(timer.Work, timer.ShortBreak)

	assert.NotEmpty(t, result)
	assert.Contains(t, result, "Work session complete!")
	assert.Contains(t, result, "Time for a short break")
}

func TestRenderComplete_ShortBreakComplete(t *testing.T) {
	result := RenderComplete(timer.ShortBreak, timer.Work)

	assert.NotEmpty(t, result)
	assert.Contains(t, result, "Break's over!")
	assert.Contains(t, result, "Ready to focus")
}

func TestRenderComplete_LongBreakComplete(t *testing.T) {
	result := RenderComplete(timer.LongBreak, timer.Work)

	assert.NotEmpty(t, result)
	assert.Contains(t, result, "Long break complete! Great work!")
	assert.Contains(t, result, "Ready to focus")
}

func TestRenderComplete_NextSessionInfo(t *testing.T) {
	tests := []struct {
		name        string
		completed   timer.SessionType
		next        timer.SessionType
		containsMsg string
	}{
		{"next is work", timer.ShortBreak, timer.Work, "Ready to focus"},
		{"next is short break", timer.Work, timer.ShortBreak, "Time for a short break"},
		{"next is long break", timer.Work, timer.LongBreak, "earned a long break"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RenderComplete(tt.completed, tt.next)
			assert.Contains(t, result, tt.containsMsg)
		})
	}
}

func TestRenderComplete_ShowsActionHint(t *testing.T) {
	result := RenderComplete(timer.Work, timer.ShortBreak)

	assert.Contains(t, result, "Press ENTER or SPACE to start")
	assert.Contains(t, result, "q to quit")
}

func TestRenderFlash(t *testing.T) {
	result := RenderFlash(80, 24)

	assert.NotEmpty(t, result)
	assert.Contains(t, result, "SESSION COMPLETE!")
}

func TestRenderFlash_DifferentSizes(t *testing.T) {
	tests := []struct {
		name   string
		width  int
		height int
	}{
		{"standard", 80, 24},
		{"wide", 120, 40},
		{"narrow", 40, 20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RenderFlash(tt.width, tt.height)

			assert.NotEmpty(t, result)
			assert.Contains(t, result, "SESSION COMPLETE!")
		})
	}
}
