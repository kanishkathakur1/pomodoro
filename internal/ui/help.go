package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// HelpItem represents a single help entry
type HelpItem struct {
	Key  string
	Desc string
}

// Default help items
var helpItems = []HelpItem{
	{"space/enter", "start/pause timer"},
	{"s", "skip session"},
	{"r", "reset timer"},
	{"n", "toggle notifications"},
	{"?", "toggle help"},
	{"q/ctrl+c", "quit"},
}

// RenderHelp creates the help overlay
func RenderHelp() string {
	var content strings.Builder

	// Title
	titleStyle := lipgloss.NewStyle().
		Foreground(Cyan).
		Bold(true).
		MarginBottom(1)

	content.WriteString(titleStyle.Render("⌨ Keyboard Shortcuts"))
	content.WriteString("\n\n")

	// Help items
	for _, item := range helpItems {
		key := HelpKeyStyle.Render(item.Key)
		desc := HelpDescStyle.Render(item.Desc)
		content.WriteString(key + desc + "\n")
	}

	// Wrap in overlay style
	return HelpOverlayStyle.Render(content.String())
}

// RenderHelpCentered renders the help overlay centered in the terminal
func RenderHelpCentered(width, height int) string {
	help := RenderHelp()
	return lipgloss.Place(
		width, height,
		lipgloss.Center, lipgloss.Center,
		help,
	)
}
