package ui

import "github.com/charmbracelet/lipgloss"

// Cyberpunk color palette
var (
	// Primary colors
	Cyan       = lipgloss.Color("#00FFFF")
	Magenta    = lipgloss.Color("#FF00FF")
	HotPink    = lipgloss.Color("#FF1493")
	ElectricBlue = lipgloss.Color("#00BFFF")
	Neon       = lipgloss.Color("#39FF14")
	Purple     = lipgloss.Color("#9D00FF")
	Yellow     = lipgloss.Color("#FFFF00")

	// Background and neutral
	DarkBg     = lipgloss.Color("#0D0D0D")
	DarkGray   = lipgloss.Color("#1A1A2E")
	MidGray    = lipgloss.Color("#333355")
	LightGray  = lipgloss.Color("#666699")

	// Session-specific colors
	WorkColor      = HotPink
	ShortBreakColor = Cyan
	LongBreakColor  = Purple
)

// Base styles
var (
	// Container style for the whole app
	AppStyle = lipgloss.NewStyle().
		Background(DarkBg).
		Padding(1, 2)

	// Title styles
	TitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(Cyan).
		MarginBottom(1)

	// ASCII timer number style
	TimerStyle = lipgloss.NewStyle().
		Foreground(HotPink).
		Bold(true)

	// Progress bar styles
	ProgressBarFilled = lipgloss.NewStyle().
		Foreground(Magenta).
		Background(Magenta)

	ProgressBarEmpty = lipgloss.NewStyle().
		Foreground(DarkGray).
		Background(DarkGray)

	// Session info style
	SessionInfoStyle = lipgloss.NewStyle().
		Foreground(Cyan).
		Bold(true).
		MarginTop(1).
		MarginBottom(1)

	// Help text style
	HelpStyle = lipgloss.NewStyle().
		Foreground(LightGray).
		MarginTop(1)

	// Key hint style
	KeyStyle = lipgloss.NewStyle().
		Foreground(Yellow).
		Bold(true)

	// Status styles
	PausedStyle = lipgloss.NewStyle().
		Foreground(Yellow).
		Bold(true).
		Blink(true)

	RunningStyle = lipgloss.NewStyle().
		Foreground(Neon).
		Bold(true)

	// Completion message style
	CompletionStyle = lipgloss.NewStyle().
		Foreground(Neon).
		Bold(true).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Cyan)

	// Splash screen styles
	SplashTitleStyle = lipgloss.NewStyle().
		Foreground(Magenta).
		Bold(true)

	SplashSubtitleStyle = lipgloss.NewStyle().
		Foreground(Cyan).
		Italic(true)

	// Help overlay styles
	HelpOverlayStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Cyan).
		Padding(1, 2).
		Background(DarkGray)

	HelpKeyStyle = lipgloss.NewStyle().
		Foreground(HotPink).
		Bold(true).
		Width(12)

	HelpDescStyle = lipgloss.NewStyle().
		Foreground(LightGray)
)

// GetSessionColor returns the appropriate color for a session type
func GetSessionColor(sessionType string) lipgloss.Color {
	switch sessionType {
	case "work":
		return WorkColor
	case "short_break":
		return ShortBreakColor
	case "long_break":
		return LongBreakColor
	default:
		return Cyan
	}
}

// GetSessionStyle returns styled text for the session type
func GetSessionStyle(sessionType string) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(GetSessionColor(sessionType)).
		Bold(true)
}
