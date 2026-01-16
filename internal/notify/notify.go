package notify

import (
	"fmt"

	"github.com/gen2brain/beeep"
	"github.com/kanishkathakur1/pomodoro/internal/config"
)

// Injectable functions for testing
var (
	systemNotifyFunc = beeep.Notify
	terminalBellFunc = func() { fmt.Print("\a") }
)

// SetNotifyFuncsForTesting replaces notification functions for testing
func SetNotifyFuncsForTesting(sn func(string, string, any) error, tb func()) {
	if sn != nil {
		systemNotifyFunc = sn
	}
	if tb != nil {
		terminalBellFunc = tb
	}
}

// ResetNotifyFuncsForTesting restores default notification functions
func ResetNotifyFuncsForTesting() {
	systemNotifyFunc = beeep.Notify
	terminalBellFunc = func() { fmt.Print("\a") }
}

// Notifier handles all notification methods
type Notifier struct {
	config *config.Config
}

// New creates a new Notifier with the given configuration
func New(cfg *config.Config) *Notifier {
	return &Notifier{config: cfg}
}

// Notify triggers all enabled notification methods
func (n *Notifier) Notify(title, message string) error {
	var lastErr error

	if n.config.Notifications.TerminalBell {
		terminalBellFunc()
	}

	if n.config.Notifications.SystemNotification {
		if err := systemNotifyFunc(title, message, ""); err != nil {
			lastErr = err
		}
	}

	return lastErr
}

// VisualFlash returns whether visual flash is enabled
func (n *Notifier) VisualFlash() bool {
	return n.config.Notifications.VisualFlash
}

// ToggleVisualFlash toggles the visual flash setting
func (n *Notifier) ToggleVisualFlash() {
	n.config.Notifications.VisualFlash = !n.config.Notifications.VisualFlash
}

// ToggleTerminalBell toggles the terminal bell setting
func (n *Notifier) ToggleTerminalBell() {
	n.config.Notifications.TerminalBell = !n.config.Notifications.TerminalBell
}

// ToggleSystemNotification toggles the system notification setting
func (n *Notifier) ToggleSystemNotification() {
	n.config.Notifications.SystemNotification = !n.config.Notifications.SystemNotification
}
