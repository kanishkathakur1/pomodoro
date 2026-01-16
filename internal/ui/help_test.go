package ui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRenderHelp_ContainsAllKeyBindings(t *testing.T) {
	result := RenderHelp()

	// Should contain all key bindings
	expectedKeys := []string{
		"space/enter",
		"start/pause",
		"s",
		"skip session",
		"r",
		"reset timer",
		"n",
		"toggle notifications",
		"?",
		"toggle help",
		"q/ctrl+c",
		"quit",
	}

	for _, key := range expectedKeys {
		assert.Contains(t, result, key, "help should contain %q", key)
	}
}

func TestRenderHelp_ContainsTitle(t *testing.T) {
	result := RenderHelp()

	assert.Contains(t, result, "Keyboard Shortcuts")
}

func TestRenderHelp_NotEmpty(t *testing.T) {
	result := RenderHelp()

	assert.NotEmpty(t, result)
}

func TestRenderHelpCentered(t *testing.T) {
	tests := []struct {
		name   string
		width  int
		height int
	}{
		{"standard 80x24", 80, 24},
		{"wide screen", 120, 40},
		{"narrow screen", 40, 20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RenderHelpCentered(tt.width, tt.height)

			// Should produce output
			assert.NotEmpty(t, result)

			// Should contain help content
			assert.Contains(t, result, "Keyboard Shortcuts")
			assert.Contains(t, result, "space/enter")
		})
	}
}

func TestRenderHelpCentered_RespectsWidth(t *testing.T) {
	// Small width should still produce valid output
	result := RenderHelpCentered(20, 10)
	assert.NotEmpty(t, result)
}

func TestRenderHelpCentered_RespectsHeight(t *testing.T) {
	// Small height should still produce valid output
	result := RenderHelpCentered(80, 5)
	assert.NotEmpty(t, result)
}

func TestHelpItems_AllDefined(t *testing.T) {
	// Verify helpItems slice has expected entries
	assert.GreaterOrEqual(t, len(helpItems), 6, "should have at least 6 help items")

	// Each item should have both key and description
	for i, item := range helpItems {
		assert.NotEmpty(t, item.Key, "help item %d should have a key", i)
		assert.NotEmpty(t, item.Desc, "help item %d should have a description", i)
	}
}
