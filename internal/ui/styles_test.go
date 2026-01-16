package ui

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/stretchr/testify/assert"
)

func TestGetSessionColor(t *testing.T) {
	tests := []struct {
		sessionType string
		expected    lipgloss.Color
	}{
		{"work", WorkColor},
		{"short_break", ShortBreakColor},
		{"long_break", LongBreakColor},
		{"unknown", Cyan},
		{"", Cyan},
	}

	for _, tt := range tests {
		t.Run(tt.sessionType, func(t *testing.T) {
			result := GetSessionColor(tt.sessionType)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetSessionStyle(t *testing.T) {
	tests := []struct {
		sessionType   string
		expectedColor lipgloss.Color
	}{
		{"work", WorkColor},
		{"short_break", ShortBreakColor},
		{"long_break", LongBreakColor},
	}

	for _, tt := range tests {
		t.Run(tt.sessionType, func(t *testing.T) {
			style := GetSessionStyle(tt.sessionType)

			// Style should render text without errors
			rendered := style.Render("test")
			assert.NotEmpty(t, rendered)
			// The rendered text should at minimum contain the original text
			assert.Contains(t, rendered, "test")
		})
	}
}

func TestColorConstants(t *testing.T) {
	// Verify color constants are defined and valid
	colors := map[string]lipgloss.Color{
		"Cyan":            Cyan,
		"Magenta":         Magenta,
		"HotPink":         HotPink,
		"ElectricBlue":    ElectricBlue,
		"Neon":            Neon,
		"Purple":          Purple,
		"Yellow":          Yellow,
		"DarkBg":          DarkBg,
		"DarkGray":        DarkGray,
		"MidGray":         MidGray,
		"LightGray":       LightGray,
		"WorkColor":       WorkColor,
		"ShortBreakColor": ShortBreakColor,
		"LongBreakColor":  LongBreakColor,
	}

	for name, color := range colors {
		t.Run(name, func(t *testing.T) {
			// Color should be a valid hex color string
			assert.NotEmpty(t, string(color), "color %s should not be empty", name)
			assert.Contains(t, string(color), "#", "color %s should be a hex color", name)
		})
	}
}

func TestStyleConstants(t *testing.T) {
	// Verify style constants can render text without errors
	styles := map[string]lipgloss.Style{
		"AppStyle":           AppStyle,
		"TitleStyle":         TitleStyle,
		"TimerStyle":         TimerStyle,
		"ProgressBarFilled":  ProgressBarFilled,
		"ProgressBarEmpty":   ProgressBarEmpty,
		"SessionInfoStyle":   SessionInfoStyle,
		"HelpStyle":          HelpStyle,
		"KeyStyle":           KeyStyle,
		"PausedStyle":        PausedStyle,
		"RunningStyle":       RunningStyle,
		"CompletionStyle":    CompletionStyle,
		"SplashTitleStyle":   SplashTitleStyle,
		"SplashSubtitleStyle": SplashSubtitleStyle,
		"HelpOverlayStyle":   HelpOverlayStyle,
		"HelpKeyStyle":       HelpKeyStyle,
		"HelpDescStyle":      HelpDescStyle,
	}

	for name, style := range styles {
		t.Run(name, func(t *testing.T) {
			// Style should render text without panicking
			result := style.Render("test")
			assert.NotEmpty(t, result, "style %s should produce output", name)
		})
	}
}

func TestWorkColorIsHotPink(t *testing.T) {
	assert.Equal(t, HotPink, WorkColor)
}

func TestShortBreakColorIsCyan(t *testing.T) {
	assert.Equal(t, Cyan, ShortBreakColor)
}

func TestLongBreakColorIsPurple(t *testing.T) {
	assert.Equal(t, Purple, LongBreakColor)
}
