package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// RenderProgressBar creates a styled progress bar
func RenderProgressBar(percent float64, width int) string {
	// Ensure percent is between 0 and 1
	if percent < 0 {
		percent = 0
	}
	if percent > 1 {
		percent = 1
	}

	// Account for percentage text at the end (e.g., " 100%")
	barWidth := width - 6
	if barWidth < 10 {
		barWidth = 10
	}

	filled := int(float64(barWidth) * percent)
	empty := barWidth - filled

	// Build the bar
	filledBar := strings.Repeat("█", filled)
	emptyBar := strings.Repeat("░", empty)

	// Style the parts
	filledStyle := lipgloss.NewStyle().Foreground(Magenta)
	emptyStyle := lipgloss.NewStyle().Foreground(DarkGray)

	bar := filledStyle.Render(filledBar) + emptyStyle.Render(emptyBar)

	// Add percentage
	percentText := fmt.Sprintf(" %3d%%", int(percent*100))
	percentStyle := lipgloss.NewStyle().Foreground(Cyan)

	return bar + percentStyle.Render(percentText)
}
