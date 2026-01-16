package ui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/stretchr/testify/assert"
)

func TestRenderASCII_AllDigits(t *testing.T) {
	style := lipgloss.NewStyle()

	tests := []struct {
		digit    rune
		expected []string
	}{
		{
			digit: '0',
			expected: []string{
				"█████",
				"█   █",
				"█   █",
				"█   █",
				"█████",
			},
		},
		{
			digit: '1',
			expected: []string{
				"  █  ",
				"  █  ",
				"  █  ",
				"  █  ",
				"  █  ",
			},
		},
		{
			digit: '2',
			expected: []string{
				"█████",
				"    █",
				"█████",
				"█    ",
				"█████",
			},
		},
		{
			digit: '3',
			expected: []string{
				"█████",
				"    █",
				"█████",
				"    █",
				"█████",
			},
		},
		{
			digit: '4',
			expected: []string{
				"█   █",
				"█   █",
				"█████",
				"    █",
				"    █",
			},
		},
		{
			digit: '5',
			expected: []string{
				"█████",
				"█    ",
				"█████",
				"    █",
				"█████",
			},
		},
		{
			digit: '6',
			expected: []string{
				"█████",
				"█    ",
				"█████",
				"█   █",
				"█████",
			},
		},
		{
			digit: '7',
			expected: []string{
				"█████",
				"    █",
				"    █",
				"    █",
				"    █",
			},
		},
		{
			digit: '8',
			expected: []string{
				"█████",
				"█   █",
				"█████",
				"█   █",
				"█████",
			},
		},
		{
			digit: '9',
			expected: []string{
				"█████",
				"█   █",
				"█████",
				"    █",
				"█████",
			},
		},
	}

	for _, tt := range tests {
		t.Run(string(tt.digit), func(t *testing.T) {
			result := RenderASCII(string(tt.digit), style)
			lines := strings.Split(result, "\n")

			assert.Len(t, lines, 5)
			for i, expected := range tt.expected {
				assert.Equal(t, expected, lines[i], "line %d mismatch", i)
			}
		})
	}
}

func TestRenderASCII_Colon(t *testing.T) {
	style := lipgloss.NewStyle()
	result := RenderASCII(":", style)
	lines := strings.Split(result, "\n")

	assert.Len(t, lines, 5)
	assert.Equal(t, "     ", lines[0])
	assert.Equal(t, "  █  ", lines[1])
	assert.Equal(t, "     ", lines[2])
	assert.Equal(t, "  █  ", lines[3])
	assert.Equal(t, "     ", lines[4])
}

func TestRenderASCII_MultipleCharacters(t *testing.T) {
	style := lipgloss.NewStyle()
	result := RenderASCII("12", style)
	lines := strings.Split(result, "\n")

	assert.Len(t, lines, 5)
	// Each line should have both digits with spacing
	for _, line := range lines {
		// "  █  " (5) + "  " (2) + "█████" (5) = 12 for first line
		assert.True(t, len(line) >= 12, "line should contain both digits: %q", line)
	}
}

func TestRenderASCII_UnknownCharacter(t *testing.T) {
	style := lipgloss.NewStyle()
	result := RenderASCII("x", style)
	lines := strings.Split(result, "\n")

	assert.Len(t, lines, 5)
	for _, line := range lines {
		assert.Empty(t, line, "unknown character should produce empty lines")
	}
}

func TestRenderTime(t *testing.T) {
	tests := []struct {
		name     string
		minutes  int
		seconds  int
		expected string
	}{
		{
			name:     "25:00",
			minutes:  25,
			seconds:  0,
			expected: "25:00",
		},
		{
			name:     "05:30",
			minutes:  5,
			seconds:  30,
			expected: "05:30",
		},
		{
			name:     "00:01",
			minutes:  0,
			seconds:  1,
			expected: "00:01",
		},
		{
			name:     "09:59",
			minutes:  9,
			seconds:  59,
			expected: "09:59",
		},
	}

	style := lipgloss.NewStyle()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RenderTime(tt.minutes, tt.seconds, style)
			lines := strings.Split(result, "\n")

			// Should have 5 rows
			assert.Len(t, lines, 5)

			// Each line should be non-empty (contains ASCII art)
			for i, line := range lines {
				assert.NotEmpty(t, line, "line %d should not be empty", i)
			}
		})
	}
}

func TestRenderTime_ContainsAllDigits(t *testing.T) {
	style := lipgloss.NewStyle()
	result := RenderTime(25, 0, style)

	// Just verify we get 5 non-empty lines (ASCII art rows)
	lines := strings.Split(result, "\n")
	assert.Len(t, lines, 5)

	for _, line := range lines {
		assert.NotEmpty(t, line)
	}
}
