package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// Config is the root structure of the config file.
type Config struct {
	Shell     string                `json:"shell"`
	ThemeName string                `json:"theme_name"`
	Segments  []string              `json:"segments"`
	Theme     map[string]ColorTheme `json:"theme"`
	Path      PathOptions           `json:"path"`
	Time      TimeOptions           `json:"time"`
}

// TimeOptions holds time-segment-specific configuration.
type TimeOptions struct {
	Format string `json:"format"` // Go time.Format string; "" = "15:04"
}

// ColorTheme holds optional fg/bg color overrides for a segment.
type ColorTheme struct {
	FG *int `json:"fg"` // nil = use segment default
	BG *int `json:"bg"`
}

// PathOptions holds path-segment-specific configuration.
type PathOptions struct {
	MaxDepth int `json:"max_depth"` // 0 = use default (4)
}

// DefaultConfigPath returns the default config file path, respecting XDG_CONFIG_HOME.
func DefaultConfigPath() string {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "highline", "config.json")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "highline", "config.json")
}

// Load reads and parses the config file at path.
// If the file does not exist, a zero-value Config is returned without error.
// If the file exists but cannot be parsed, an error is returned.
func Load(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Config{}, nil
		}
		return Config{}, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("parse config: %w", err)
	}

	// Validate and sanitize color values (must be 0–15).
	for name, theme := range cfg.Theme {
		modified := false
		if theme.FG != nil && (*theme.FG < 0 || *theme.FG > 15) {
			fmt.Fprintf(os.Stderr, "highline: config: theme.%s.fg=%d out of range [0,15], using default\n", name, *theme.FG)
			theme.FG = nil
			modified = true
		}
		if theme.BG != nil && (*theme.BG < 0 || *theme.BG > 15) {
			fmt.Fprintf(os.Stderr, "highline: config: theme.%s.bg=%d out of range [0,15], using default\n", name, *theme.BG)
			theme.BG = nil
			modified = true
		}
		if modified {
			cfg.Theme[name] = theme
		}
	}

	return cfg, nil
}
