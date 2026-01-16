package app

import "github.com/charmbracelet/bubbles/key"

// KeyMap defines all keyboard bindings
type KeyMap struct {
	Toggle key.Binding
	Skip   key.Binding
	Reset  key.Binding
	Notify key.Binding
	Help   key.Binding
	Quit   key.Binding
}

// DefaultKeyMap returns the default key bindings
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Toggle: key.NewBinding(
			key.WithKeys(" ", "enter"),
			key.WithHelp("space/enter", "start/pause"),
		),
		Skip: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "skip session"),
		),
		Reset: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "reset timer"),
		),
		Notify: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "toggle notifications"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
	}
}
