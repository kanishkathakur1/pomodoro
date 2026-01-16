package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestConfig(t *testing.T) (string, func()) {
	t.Helper()

	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "pomodoro-config-test-*")
	require.NoError(t, err)

	configFile := filepath.Join(tmpDir, "config.toml")
	SetConfigPathForTesting(configFile)

	cleanup := func() {
		ResetConfigPathForTesting()
		os.RemoveAll(tmpDir)
	}

	return configFile, cleanup
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	assert.True(t, cfg.Notifications.VisualFlash, "VisualFlash should be enabled by default")
	assert.True(t, cfg.Notifications.TerminalBell, "TerminalBell should be enabled by default")
	assert.True(t, cfg.Notifications.SystemNotification, "SystemNotification should be enabled by default")
}

func TestLoad_NoConfigFile(t *testing.T) {
	configFile, cleanup := setupTestConfig(t)
	defer cleanup()

	// Ensure file doesn't exist
	os.Remove(configFile)

	cfg, err := Load()

	require.NoError(t, err)
	// Should return defaults
	assert.True(t, cfg.Notifications.VisualFlash)
	assert.True(t, cfg.Notifications.TerminalBell)
	assert.True(t, cfg.Notifications.SystemNotification)

	// File should now exist (created with defaults)
	_, err = os.Stat(configFile)
	assert.NoError(t, err, "config file should be created")
}

func TestLoad_ValidConfigFile(t *testing.T) {
	configFile, cleanup := setupTestConfig(t)
	defer cleanup()

	// Create config file with custom values
	configContent := `[notifications]
visual_flash = false
terminal_bell = true
system_notification = false
`
	err := os.MkdirAll(filepath.Dir(configFile), 0755)
	require.NoError(t, err)
	err = os.WriteFile(configFile, []byte(configContent), 0644)
	require.NoError(t, err)

	cfg, err := Load()

	require.NoError(t, err)
	assert.False(t, cfg.Notifications.VisualFlash)
	assert.True(t, cfg.Notifications.TerminalBell)
	assert.False(t, cfg.Notifications.SystemNotification)
}

func TestLoad_InvalidTOML(t *testing.T) {
	configFile, cleanup := setupTestConfig(t)
	defer cleanup()

	// Create invalid TOML file
	err := os.MkdirAll(filepath.Dir(configFile), 0755)
	require.NoError(t, err)
	err = os.WriteFile(configFile, []byte("this is { not valid toml"), 0644)
	require.NoError(t, err)

	cfg, err := Load()

	// Should return defaults on parse error
	require.NoError(t, err)
	assert.True(t, cfg.Notifications.VisualFlash)
	assert.True(t, cfg.Notifications.TerminalBell)
	assert.True(t, cfg.Notifications.SystemNotification)
}

func TestSave(t *testing.T) {
	configFile, cleanup := setupTestConfig(t)
	defer cleanup()

	cfg := &Config{
		Notifications: NotificationConfig{
			VisualFlash:        false,
			TerminalBell:       true,
			SystemNotification: false,
		},
	}

	err := cfg.Save()
	require.NoError(t, err)

	// Verify file exists and has correct content
	content, err := os.ReadFile(configFile)
	require.NoError(t, err)

	assert.Contains(t, string(content), "visual_flash = false")
	assert.Contains(t, string(content), "terminal_bell = true")
	assert.Contains(t, string(content), "system_notification = false")
}

func TestSave_CreatesDirectory(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "pomodoro-config-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Set config path to a nested directory that doesn't exist
	nestedPath := filepath.Join(tmpDir, "subdir", "another", "config.toml")
	SetConfigPathForTesting(nestedPath)
	defer ResetConfigPathForTesting()

	cfg := DefaultConfig()
	err = cfg.Save()

	require.NoError(t, err)
	_, err = os.Stat(nestedPath)
	assert.NoError(t, err, "config file should be created in nested directory")
}

func TestLoadSaveRoundTrip(t *testing.T) {
	configFile, cleanup := setupTestConfig(t)
	defer cleanup()

	// Ensure config file doesn't exist
	os.Remove(configFile)

	// Create custom config
	original := &Config{
		Notifications: NotificationConfig{
			VisualFlash:        false,
			TerminalBell:       false,
			SystemNotification: true,
		},
	}

	// Save it
	err := original.Save()
	require.NoError(t, err)

	// Load it back
	loaded, err := Load()
	require.NoError(t, err)

	// Compare
	assert.Equal(t, original.Notifications.VisualFlash, loaded.Notifications.VisualFlash)
	assert.Equal(t, original.Notifications.TerminalBell, loaded.Notifications.TerminalBell)
	assert.Equal(t, original.Notifications.SystemNotification, loaded.Notifications.SystemNotification)
}

func TestLoad_PartialConfig(t *testing.T) {
	configFile, cleanup := setupTestConfig(t)
	defer cleanup()

	// Create config file with only some values (missing terminal_bell)
	configContent := `[notifications]
visual_flash = false
system_notification = true
`
	err := os.MkdirAll(filepath.Dir(configFile), 0755)
	require.NoError(t, err)
	err = os.WriteFile(configFile, []byte(configContent), 0644)
	require.NoError(t, err)

	cfg, err := Load()

	require.NoError(t, err)
	assert.False(t, cfg.Notifications.VisualFlash)
	// Missing field should be zero value (false for bool)
	assert.False(t, cfg.Notifications.TerminalBell)
	assert.True(t, cfg.Notifications.SystemNotification)
}

func TestConfigPath_UsesOverride(t *testing.T) {
	customPath := "/custom/path/config.toml"
	SetConfigPathForTesting(customPath)
	defer ResetConfigPathForTesting()

	path, err := configPath()

	require.NoError(t, err)
	assert.Equal(t, customPath, path)
}

func TestConfigPath_UsesDefaultWhenNoOverride(t *testing.T) {
	ResetConfigPathForTesting()

	path, err := configPath()

	require.NoError(t, err)
	assert.Contains(t, path, "pomodoro")
	assert.Contains(t, path, "config.toml")
}
