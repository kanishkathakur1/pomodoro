package app

import (
	"testing"

	"github.com/charmbracelet/bubbles/key"
	"github.com/stretchr/testify/assert"
)

func TestDefaultKeyMap(t *testing.T) {
	km := DefaultKeyMap()

	// Verify all bindings are defined
	assert.NotEmpty(t, km.Toggle.Keys(), "Toggle should have keys")
	assert.NotEmpty(t, km.Skip.Keys(), "Skip should have keys")
	assert.NotEmpty(t, km.Reset.Keys(), "Reset should have keys")
	assert.NotEmpty(t, km.Notify.Keys(), "Notify should have keys")
	assert.NotEmpty(t, km.Help.Keys(), "Help should have keys")
	assert.NotEmpty(t, km.Quit.Keys(), "Quit should have keys")
}

func TestDefaultKeyMap_ToggleKeys(t *testing.T) {
	km := DefaultKeyMap()

	// Toggle should be space and enter
	assert.Contains(t, km.Toggle.Keys(), " ")
	assert.Contains(t, km.Toggle.Keys(), "enter")
}

func TestDefaultKeyMap_SkipKey(t *testing.T) {
	km := DefaultKeyMap()

	// Skip should be 's'
	assert.Contains(t, km.Skip.Keys(), "s")
}

func TestDefaultKeyMap_ResetKey(t *testing.T) {
	km := DefaultKeyMap()

	// Reset should be 'r'
	assert.Contains(t, km.Reset.Keys(), "r")
}

func TestDefaultKeyMap_NotifyKey(t *testing.T) {
	km := DefaultKeyMap()

	// Notify should be 'n'
	assert.Contains(t, km.Notify.Keys(), "n")
}

func TestDefaultKeyMap_HelpKey(t *testing.T) {
	km := DefaultKeyMap()

	// Help should be '?'
	assert.Contains(t, km.Help.Keys(), "?")
}

func TestDefaultKeyMap_QuitKeys(t *testing.T) {
	km := DefaultKeyMap()

	// Quit should be 'q' and 'ctrl+c'
	assert.Contains(t, km.Quit.Keys(), "q")
	assert.Contains(t, km.Quit.Keys(), "ctrl+c")
}

func TestDefaultKeyMap_HelpText(t *testing.T) {
	km := DefaultKeyMap()

	// Verify help text is set
	toggleHelp := km.Toggle.Help()
	assert.NotEmpty(t, toggleHelp.Key)
	assert.NotEmpty(t, toggleHelp.Desc)
	assert.Equal(t, "space/enter", toggleHelp.Key)
	assert.Equal(t, "start/pause", toggleHelp.Desc)

	skipHelp := km.Skip.Help()
	assert.Equal(t, "s", skipHelp.Key)
	assert.Equal(t, "skip session", skipHelp.Desc)

	resetHelp := km.Reset.Help()
	assert.Equal(t, "r", resetHelp.Key)
	assert.Equal(t, "reset timer", resetHelp.Desc)

	notifyHelp := km.Notify.Help()
	assert.Equal(t, "n", notifyHelp.Key)
	assert.Equal(t, "toggle notifications", notifyHelp.Desc)

	helpHelp := km.Help.Help()
	assert.Equal(t, "?", helpHelp.Key)
	assert.Equal(t, "toggle help", helpHelp.Desc)

	quitHelp := km.Quit.Help()
	assert.Equal(t, "q", quitHelp.Key)
	assert.Equal(t, "quit", quitHelp.Desc)
}

func TestKeyMap_AllBindingsHaveHelp(t *testing.T) {
	km := DefaultKeyMap()

	bindings := []key.Binding{
		km.Toggle,
		km.Skip,
		km.Reset,
		km.Notify,
		km.Help,
		km.Quit,
	}

	for _, binding := range bindings {
		help := binding.Help()
		assert.NotEmpty(t, help.Key, "binding should have help key")
		assert.NotEmpty(t, help.Desc, "binding should have help description")
	}
}
