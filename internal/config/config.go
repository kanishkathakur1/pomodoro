package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Config holds all application configuration
type Config struct {
	Notifications NotificationConfig `toml:"notifications"`
}

// NotificationConfig controls notification behavior
type NotificationConfig struct {
	VisualFlash        bool `toml:"visual_flash"`
	TerminalBell       bool `toml:"terminal_bell"`
	SystemNotification bool `toml:"system_notification"`
}

// configPathOverride allows tests to inject a custom config path
var configPathOverride string

// SetConfigPathForTesting sets a custom config path for testing purposes
func SetConfigPathForTesting(path string) {
	configPathOverride = path
}

// ResetConfigPathForTesting resets the config path override
func ResetConfigPathForTesting() {
	configPathOverride = ""
}

// DefaultConfig returns sensible default configuration
func DefaultConfig() *Config {
	return &Config{
		Notifications: NotificationConfig{
			VisualFlash:        true,
			TerminalBell:       true,
			SystemNotification: true,
		},
	}
}

// configPath returns the path to the config file
func configPath() (string, error) {
	if configPathOverride != "" {
		return configPathOverride, nil
	}
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "pomodoro", "config.toml"), nil
}

// Load reads configuration from the config file
// If the file doesn't exist, it creates one with defaults
func Load() (*Config, error) {
	path, err := configPath()
	if err != nil {
		return DefaultConfig(), nil
	}

	// Check if config file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Create default config
		cfg := DefaultConfig()
		if err := cfg.Save(); err != nil {
			// If we can't save, just return defaults
			return cfg, nil
		}
		return cfg, nil
	}

	// Read existing config
	cfg := &Config{}
	if _, err := toml.DecodeFile(path, cfg); err != nil {
		// On error, return defaults
		return DefaultConfig(), nil
	}

	return cfg, nil
}

// Save writes the configuration to the config file
func (c *Config) Save() error {
	path, err := configPath()
	if err != nil {
		return err
	}

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Write config file
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := toml.NewEncoder(f)
	return encoder.Encode(c)
}
