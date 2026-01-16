package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/kanishkathakur1/pomodoro/internal/timer"
)

// Splash screen ASCII art
var splashTitle = []string{
	"██████╗  ██████╗ ███╗   ███╗ ██████╗ ██████╗  ██████╗ ██████╗  ██████╗ ",
	"██╔══██╗██╔═══██╗████╗ ████║██╔═══██╗██╔══██╗██╔═══██╗██╔══██╗██╔═══██╗",
	"██████╔╝██║   ██║██╔████╔██║██║   ██║██║  ██║██║   ██║██████╔╝██║   ██║",
	"██╔═══╝ ██║   ██║██║╚██╔╝██║██║   ██║██║  ██║██║   ██║██╔══██╗██║   ██║",
	"██║     ╚██████╔╝██║ ╚═╝ ██║╚██████╔╝██████╔╝╚██████╔╝██║  ██║╚██████╔╝",
	"╚═╝      ╚═════╝ ╚═╝     ╚═╝ ╚═════╝ ╚═════╝  ╚═════╝ ╚═╝  ╚═╝ ╚═════╝ ",
}

// RenderSplash creates the animated splash screen
func RenderSplash(frame int, width, height int) string {
	// Animate colors based on frame
	colors := []lipgloss.Color{Magenta, HotPink, Cyan, ElectricBlue, Purple}
	colorIndex := frame % len(colors)
	titleColor := colors[colorIndex]

	titleStyle := lipgloss.NewStyle().Foreground(titleColor).Bold(true)
	subtitleStyle := SplashSubtitleStyle

	// Build the splash content
	var content strings.Builder

	for _, line := range splashTitle {
		content.WriteString(titleStyle.Render(line))
		content.WriteString("\n")
	}

	content.WriteString("\n")
	content.WriteString(subtitleStyle.Render("Focus. Flow. Flourish."))
	content.WriteString("\n\n")
	content.WriteString(subtitleStyle.Render("Press any key to start..."))

	// Center the content
	return lipgloss.Place(
		width, height,
		lipgloss.Center, lipgloss.Center,
		content.String(),
	)
}

// RenderTimer renders the main timer view
func RenderTimer(t *timer.Timer, width, height int, paused bool) string {
	// Get session-appropriate color
	sessionColor := GetSessionColor(string(t.SessionType))
	timerStyle := lipgloss.NewStyle().Foreground(sessionColor).Bold(true)

	var content strings.Builder

	// Session title
	titleStyle := lipgloss.NewStyle().Foreground(sessionColor).Bold(true)
	content.WriteString(titleStyle.Render(t.SessionName()))
	content.WriteString("\n\n")

	// ASCII time display
	content.WriteString(RenderTime(t.MinutesRemaining(), t.SecondsRemaining(), timerStyle))
	content.WriteString("\n\n")

	// Progress bar
	progressWidth := 50
	if width < 60 {
		progressWidth = width - 10
	}
	content.WriteString(RenderProgressBar(t.Progress(), progressWidth))
	content.WriteString("\n\n")

	// Session counter
	sessionInfo := fmt.Sprintf("Pomodoro %d/%d", t.PomodoroCount, timer.PomodorosBeforeLongBreak)
	if t.SessionType == timer.Work {
		if t.PomodoroCount >= timer.PomodorosBeforeLongBreak {
			sessionInfo += " • Long break next!"
		} else {
			sessionInfo += " • Short break next"
		}
	} else {
		sessionInfo += " • Work session next"
	}
	content.WriteString(SessionInfoStyle.Render(sessionInfo))
	content.WriteString("\n")

	// Status indicator
	if paused {
		content.WriteString(PausedStyle.Render("⏸ PAUSED"))
	} else {
		content.WriteString(RunningStyle.Render("▶ RUNNING"))
	}
	content.WriteString("\n\n")

	// Help hint
	content.WriteString(HelpStyle.Render("Press ? for help • q to quit"))

	// Center the content
	return lipgloss.Place(
		width, height,
		lipgloss.Center, lipgloss.Center,
		content.String(),
	)
}

// RenderComplete renders the session complete view
func RenderComplete(completedSession timer.SessionType, nextSession timer.SessionType) string {
	var content strings.Builder

	// Completion message
	var completeMsg string
	switch completedSession {
	case timer.Work:
		completeMsg = "🎉 Work session complete!"
	case timer.ShortBreak:
		completeMsg = "☕ Break's over!"
	case timer.LongBreak:
		completeMsg = "🌟 Long break complete! Great work!"
	}

	content.WriteString(CompletionStyle.Render(completeMsg))
	content.WriteString("\n\n")

	// Next session info
	var nextMsg string
	switch nextSession {
	case timer.Work:
		nextMsg = "Ready to focus? Start your work session."
	case timer.ShortBreak:
		nextMsg = "Time for a short break. Rest your eyes!"
	case timer.LongBreak:
		nextMsg = "You've earned a long break! Take 15 minutes."
	}
	content.WriteString(SessionInfoStyle.Render(nextMsg))
	content.WriteString("\n\n")

	// Action hint
	content.WriteString(HelpStyle.Render("Press ENTER or SPACE to start • q to quit"))

	return content.String()
}

// RenderFlash renders a visual flash effect
func RenderFlash(width, height int) string {
	flashStyle := lipgloss.NewStyle().
		Background(Cyan).
		Foreground(DarkBg).
		Bold(true)

	flashMsg := flashStyle.Render("  SESSION COMPLETE!  ")

	return lipgloss.Place(
		width, height,
		lipgloss.Center, lipgloss.Center,
		flashMsg,
	)
}
