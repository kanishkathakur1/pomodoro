package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ASCII digit representations (5 rows each)
var digits = map[rune][]string{
	'0': {
		"█████",
		"█   █",
		"█   █",
		"█   █",
		"█████",
	},
	'1': {
		"  █  ",
		"  █  ",
		"  █  ",
		"  █  ",
		"  █  ",
	},
	'2': {
		"█████",
		"    █",
		"█████",
		"█    ",
		"█████",
	},
	'3': {
		"█████",
		"    █",
		"█████",
		"    █",
		"█████",
	},
	'4': {
		"█   █",
		"█   █",
		"█████",
		"    █",
		"    █",
	},
	'5': {
		"█████",
		"█    ",
		"█████",
		"    █",
		"█████",
	},
	'6': {
		"█████",
		"█    ",
		"█████",
		"█   █",
		"█████",
	},
	'7': {
		"█████",
		"    █",
		"    █",
		"    █",
		"    █",
	},
	'8': {
		"█████",
		"█   █",
		"█████",
		"█   █",
		"█████",
	},
	'9': {
		"█████",
		"█   █",
		"█████",
		"    █",
		"█████",
	},
	':': {
		"     ",
		"  █  ",
		"     ",
		"  █  ",
		"     ",
	},
}

// RenderASCII renders text as ASCII art using the provided style
func RenderASCII(text string, style lipgloss.Style) string {
	rows := make([]string, 5)

	for _, char := range text {
		digit, ok := digits[char]
		if !ok {
			// Skip unknown characters
			continue
		}
		for i := 0; i < 5; i++ {
			if rows[i] != "" {
				rows[i] += "  " // spacing between characters
			}
			rows[i] += digit[i]
		}
	}

	// Apply style to each row and join
	styledRows := make([]string, 5)
	for i, row := range rows {
		styledRows[i] = style.Render(row)
	}

	return strings.Join(styledRows, "\n")
}

// RenderTime renders minutes and seconds as MM:SS ASCII art
func RenderTime(minutes, seconds int, style lipgloss.Style) string {
	// Format as MM:SS
	timeStr := ""
	timeStr += string(rune('0' + minutes/10))
	timeStr += string(rune('0' + minutes%10))
	timeStr += ":"
	timeStr += string(rune('0' + seconds/10))
	timeStr += string(rune('0' + seconds%10))

	return RenderASCII(timeStr, style)
}
