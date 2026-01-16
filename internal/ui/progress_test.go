package ui

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRenderProgressBar_Percentages(t *testing.T) {
	tests := []struct {
		name    string
		percent float64
	}{
		{"0%", 0.0},
		{"25%", 0.25},
		{"50%", 0.5},
		{"75%", 0.75},
		{"100%", 1.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RenderProgressBar(tt.percent, 50)

			// Should contain percentage text
			expectedPercent := int(tt.percent * 100)
			assert.Contains(t, result, "%")

			// Result should not be empty
			assert.NotEmpty(t, result)

			// Percentage should be shown
			if expectedPercent == 0 {
				assert.Contains(t, result, "  0%")
			} else if expectedPercent == 100 {
				assert.Contains(t, result, "100%")
			}
		})
	}
}

func TestRenderProgressBar_Clamping(t *testing.T) {
	tests := []struct {
		name    string
		percent float64
	}{
		{"negative clamped to 0", -0.5},
		{"greater than 1 clamped to 100", 1.5},
		{"much greater than 1 clamped", 10.0},
		{"much less than 0 clamped", -10.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RenderProgressBar(tt.percent, 50)

			// Should not cause errors and should produce valid output
			assert.NotEmpty(t, result)
			assert.Contains(t, result, "%")
		})
	}
}

func TestRenderProgressBar_NegativeClampedToZero(t *testing.T) {
	result := RenderProgressBar(-0.5, 50)
	// Should show 0%
	assert.Contains(t, result, "0%")
}

func TestRenderProgressBar_OverOneClampedTo100(t *testing.T) {
	result := RenderProgressBar(1.5, 50)
	// Should show 100%
	assert.Contains(t, result, "100%")
}

func TestRenderProgressBar_Widths(t *testing.T) {
	tests := []struct {
		name  string
		width int
	}{
		{"narrow width", 20},
		{"medium width", 50},
		{"wide width", 100},
		{"very narrow width", 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RenderProgressBar(0.5, tt.width)

			// Should produce valid output
			assert.NotEmpty(t, result)
			assert.Contains(t, result, "%")
		})
	}
}

func TestRenderProgressBar_MinimumBarWidth(t *testing.T) {
	// Very narrow width should still work with minimum bar width of 10
	result := RenderProgressBar(0.5, 5)

	// Should produce valid output
	assert.NotEmpty(t, result)
	// Should contain bar characters
	assert.True(t, strings.Contains(result, "█") || strings.Contains(result, "░"))
}

func TestRenderProgressBar_FilledAndEmpty(t *testing.T) {
	// At 50%, should have both filled and empty characters
	result := RenderProgressBar(0.5, 50)

	assert.Contains(t, result, "█", "should contain filled bar characters")
	assert.Contains(t, result, "░", "should contain empty bar characters")
}

func TestRenderProgressBar_FullBar(t *testing.T) {
	result := RenderProgressBar(1.0, 50)

	// At 100%, should have mostly filled characters
	assert.Contains(t, result, "█")
	assert.Contains(t, result, "100%")
}

func TestRenderProgressBar_EmptyBar(t *testing.T) {
	result := RenderProgressBar(0.0, 50)

	// At 0%, should have mostly empty characters
	assert.Contains(t, result, "░")
	assert.Contains(t, result, "0%")
}
