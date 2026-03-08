package config

import (
	"os"
	"path/filepath"
	"testing"
)

func writeConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestLoad_FileNotExist(t *testing.T) {
	cfg, err := Load("/nonexistent/path/config.json")
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if cfg.Shell != "" || len(cfg.Segments) != 0 {
		t.Errorf("expected zero-value Config, got: %+v", cfg)
	}
}

func TestLoad_ValidConfig(t *testing.T) {
	content := `{
		"shell": "zsh",
		"segments": ["git", "path"],
		"theme": {
			"path": {"fg": 7, "bg": 4},
			"git":  {"fg": 0, "bg": 2}
		},
		"path": {"max_depth": 3}
	}`
	path := writeConfig(t, content)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Shell != "zsh" {
		t.Errorf("shell: want 'zsh', got %q", cfg.Shell)
	}
	if len(cfg.Segments) != 2 || cfg.Segments[0] != "git" || cfg.Segments[1] != "path" {
		t.Errorf("segments: want [git path], got %v", cfg.Segments)
	}
	if cfg.Path.MaxDepth != 3 {
		t.Errorf("path.max_depth: want 3, got %d", cfg.Path.MaxDepth)
	}
	pathTheme := cfg.Theme["path"]
	if pathTheme.FG == nil || *pathTheme.FG != 7 {
		t.Errorf("theme.path.fg: want 7, got %v", pathTheme.FG)
	}
	if pathTheme.BG == nil || *pathTheme.BG != 4 {
		t.Errorf("theme.path.bg: want 4, got %v", pathTheme.BG)
	}
}

func TestLoad_TimeOptions(t *testing.T) {
	content := `{"time": {"format": "2006-01-02"}}`
	path := writeConfig(t, content)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Time.Format != "2006-01-02" {
		t.Errorf("time.format: want %q, got %q", "2006-01-02", cfg.Time.Format)
	}
}

func TestLoad_ThemeName(t *testing.T) {
	content := `{"theme_name": "nord"}`
	path := writeConfig(t, content)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.ThemeName != "nord" {
		t.Errorf("theme_name: want %q, got %q", "nord", cfg.ThemeName)
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	path := writeConfig(t, `{not valid json}`)
	_, err := Load(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestLoad_OutOfRangeColors(t *testing.T) {
	content := `{"theme": {"path": {"fg": 99, "bg": -1}}}`
	path := writeConfig(t, content)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	theme := cfg.Theme["path"]
	if theme.FG != nil {
		t.Errorf("expected fg to be nil (out-of-range fallback), got %v", theme.FG)
	}
	if theme.BG != nil {
		t.Errorf("expected bg to be nil (out-of-range fallback), got %v", theme.BG)
	}
}

func TestLoad_ValidBoundaryColors(t *testing.T) {
	content := `{"theme": {"path": {"fg": 0, "bg": 15}}}`
	path := writeConfig(t, content)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	theme := cfg.Theme["path"]
	if theme.FG == nil || *theme.FG != 0 {
		t.Errorf("expected fg=0, got %v", theme.FG)
	}
	if theme.BG == nil || *theme.BG != 15 {
		t.Errorf("expected bg=15, got %v", theme.BG)
	}
}

func TestDefaultConfigPath_XDG(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", "/custom/xdg")
	got := DefaultConfigPath()
	want := "/custom/xdg/highline/config.json"
	if got != want {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestDefaultConfigPath_Home(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", "")
	got := DefaultConfigPath()
	if got == "" {
		t.Error("expected non-empty path")
	}
	// Should end with the expected suffix.
	want := filepath.Join("highline", "config.json")
	if len(got) < len(want) || got[len(got)-len(want):] != want {
		t.Errorf("expected path to end with %q, got %q", want, got)
	}
}
