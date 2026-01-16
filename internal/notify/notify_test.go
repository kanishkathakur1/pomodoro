package notify

import (
	"errors"
	"testing"

	"github.com/kanishkathakur1/pomodoro/internal/config"
	"github.com/stretchr/testify/assert"
)

func setupMocks(t *testing.T) (bellCalled *bool, notifyCalled *bool, notifyTitle *string, notifyMessage *string, cleanup func()) {
	t.Helper()

	bc := false
	nc := false
	var nt, nm string

	bellCalled = &bc
	notifyCalled = &nc
	notifyTitle = &nt
	notifyMessage = &nm

	SetNotifyFuncsForTesting(
		func(title, message string, icon any) error {
			*notifyCalled = true
			*notifyTitle = title
			*notifyMessage = message
			return nil
		},
		func() {
			*bellCalled = true
		},
	)

	cleanup = func() {
		ResetNotifyFuncsForTesting()
	}

	return
}

func TestNew(t *testing.T) {
	cfg := config.DefaultConfig()
	notifier := New(cfg)

	assert.NotNil(t, notifier)
	assert.Equal(t, cfg, notifier.config)
}

func TestVisualFlash(t *testing.T) {
	tests := []struct {
		name     string
		enabled  bool
		expected bool
	}{
		{"enabled", true, true},
		{"disabled", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				Notifications: config.NotificationConfig{
					VisualFlash: tt.enabled,
				},
			}
			notifier := New(cfg)

			assert.Equal(t, tt.expected, notifier.VisualFlash())
		})
	}
}

func TestToggleVisualFlash(t *testing.T) {
	cfg := &config.Config{
		Notifications: config.NotificationConfig{
			VisualFlash: true,
		},
	}
	notifier := New(cfg)

	assert.True(t, notifier.VisualFlash())

	notifier.ToggleVisualFlash()
	assert.False(t, notifier.VisualFlash())

	notifier.ToggleVisualFlash()
	assert.True(t, notifier.VisualFlash())
}

func TestToggleTerminalBell(t *testing.T) {
	cfg := &config.Config{
		Notifications: config.NotificationConfig{
			TerminalBell: true,
		},
	}
	notifier := New(cfg)

	assert.True(t, cfg.Notifications.TerminalBell)

	notifier.ToggleTerminalBell()
	assert.False(t, cfg.Notifications.TerminalBell)

	notifier.ToggleTerminalBell()
	assert.True(t, cfg.Notifications.TerminalBell)
}

func TestToggleSystemNotification(t *testing.T) {
	cfg := &config.Config{
		Notifications: config.NotificationConfig{
			SystemNotification: true,
		},
	}
	notifier := New(cfg)

	assert.True(t, cfg.Notifications.SystemNotification)

	notifier.ToggleSystemNotification()
	assert.False(t, cfg.Notifications.SystemNotification)

	notifier.ToggleSystemNotification()
	assert.True(t, cfg.Notifications.SystemNotification)
}

func TestNotify_AllEnabled(t *testing.T) {
	bellCalled, notifyCalled, notifyTitle, notifyMessage, cleanup := setupMocks(t)
	defer cleanup()

	cfg := &config.Config{
		Notifications: config.NotificationConfig{
			VisualFlash:        true,
			TerminalBell:       true,
			SystemNotification: true,
		},
	}
	notifier := New(cfg)

	err := notifier.Notify("Test Title", "Test Message")

	assert.NoError(t, err)
	assert.True(t, *bellCalled, "terminal bell should be called")
	assert.True(t, *notifyCalled, "system notification should be called")
	assert.Equal(t, "Test Title", *notifyTitle)
	assert.Equal(t, "Test Message", *notifyMessage)
}

func TestNotify_AllDisabled(t *testing.T) {
	bellCalled, notifyCalled, _, _, cleanup := setupMocks(t)
	defer cleanup()

	cfg := &config.Config{
		Notifications: config.NotificationConfig{
			VisualFlash:        false,
			TerminalBell:       false,
			SystemNotification: false,
		},
	}
	notifier := New(cfg)

	err := notifier.Notify("Test Title", "Test Message")

	assert.NoError(t, err)
	assert.False(t, *bellCalled, "terminal bell should not be called")
	assert.False(t, *notifyCalled, "system notification should not be called")
}

func TestNotify_OnlyBell(t *testing.T) {
	bellCalled, notifyCalled, _, _, cleanup := setupMocks(t)
	defer cleanup()

	cfg := &config.Config{
		Notifications: config.NotificationConfig{
			VisualFlash:        false,
			TerminalBell:       true,
			SystemNotification: false,
		},
	}
	notifier := New(cfg)

	err := notifier.Notify("Test Title", "Test Message")

	assert.NoError(t, err)
	assert.True(t, *bellCalled, "terminal bell should be called")
	assert.False(t, *notifyCalled, "system notification should not be called")
}

func TestNotify_OnlySystemNotification(t *testing.T) {
	bellCalled, notifyCalled, _, _, cleanup := setupMocks(t)
	defer cleanup()

	cfg := &config.Config{
		Notifications: config.NotificationConfig{
			VisualFlash:        false,
			TerminalBell:       false,
			SystemNotification: true,
		},
	}
	notifier := New(cfg)

	err := notifier.Notify("Test Title", "Test Message")

	assert.NoError(t, err)
	assert.False(t, *bellCalled, "terminal bell should not be called")
	assert.True(t, *notifyCalled, "system notification should be called")
}

func TestNotify_SystemNotificationError(t *testing.T) {
	defer ResetNotifyFuncsForTesting()

	expectedErr := errors.New("notification failed")
	SetNotifyFuncsForTesting(
		func(title, message string, icon any) error {
			return expectedErr
		},
		func() {},
	)

	cfg := &config.Config{
		Notifications: config.NotificationConfig{
			VisualFlash:        false,
			TerminalBell:       false,
			SystemNotification: true,
		},
	}
	notifier := New(cfg)

	err := notifier.Notify("Test Title", "Test Message")

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestNotify_BellStillCallsOnNotificationError(t *testing.T) {
	defer ResetNotifyFuncsForTesting()

	bellCalled := false
	expectedErr := errors.New("notification failed")
	SetNotifyFuncsForTesting(
		func(title, message string, icon any) error {
			return expectedErr
		},
		func() {
			bellCalled = true
		},
	)

	cfg := &config.Config{
		Notifications: config.NotificationConfig{
			VisualFlash:        false,
			TerminalBell:       true,
			SystemNotification: true,
		},
	}
	notifier := New(cfg)

	err := notifier.Notify("Test Title", "Test Message")

	assert.Error(t, err)
	assert.True(t, bellCalled, "terminal bell should still be called even on notification error")
}

func TestSetNotifyFuncsForTesting_NilDoesNotReplace(t *testing.T) {
	defer ResetNotifyFuncsForTesting()

	called := false
	SetNotifyFuncsForTesting(nil, func() { called = true })

	terminalBellFunc()
	assert.True(t, called)
}
