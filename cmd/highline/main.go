package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/user/highline/internal/config"
	"github.com/user/highline/internal/renderer"
	"github.com/user/highline/internal/segment"
	"github.com/user/highline/internal/theme"
)

func main() {
	shell := flag.String("shell", "", "target shell: bash or zsh")
	configPath := flag.String("config", "", "path to config file")
	themeName := flag.String("theme", "", "built-in color theme (e.g. monokai, nord)")
	flag.Parse()

	// Determine config file path.
	cfgPath := *configPath
	if cfgPath == "" {
		cfgPath = config.DefaultConfigPath()
	}

	// Load config; on parse error warn and continue with defaults.
	cfg, err := config.Load(cfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "highline: config error: %v, using defaults\n", err)
	}

	// CLI flags override config values.
	if *shell != "" {
		cfg.Shell = *shell
	}
	if cfg.Shell == "" {
		cfg.Shell = "bash"
	}

	// Resolve theme name: --theme flag > config.ThemeName.
	resolvedThemeName := cfg.ThemeName
	if *themeName != "" {
		resolvedThemeName = *themeName
	}

	// Load built-in theme if specified.
	var activeTheme theme.Theme
	if resolvedThemeName != "" {
		t, ok := theme.Get(resolvedThemeName)
		if !ok {
			fmt.Fprintf(os.Stderr, "highline: unknown theme %q; available themes: %s\n",
				resolvedThemeName, strings.Join(theme.Names(), ", "))
		} else {
			activeTheme = t
		}
	}

	names := flag.Args()
	if len(names) > 0 {
		cfg.Segments = names
	}
	if len(cfg.Segments) == 0 {
		cfg.Segments = []string{"path", "git", "kube"}
	}

	// Build per-segment configs, applying built-in theme then per-segment overrides.
	pathFG, pathBG := resolveColors(activeTheme, cfg.Theme["path"], "path")
	pathCfg := segment.PathConfig{MaxDepth: cfg.Path.MaxDepth, FG: pathFG, BG: pathBG}

	gitFG, gitBG := resolveColors(activeTheme, cfg.Theme["git"], "git")
	gitCfg := segment.GitConfig{FG: gitFG, BG: gitBG}

	kubeFG, kubeBG := resolveColors(activeTheme, cfg.Theme["kube"], "kube")
	kubeCfg := segment.KubeConfig{FG: kubeFG, BG: kubeBG}

	promptFG, promptBG := resolveColors(activeTheme, cfg.Theme["prompt"], "prompt")
	promptCfg := segment.PromptConfig{Shell: cfg.Shell, FG: promptFG, BG: promptBG}

	timeFG, timeBG := resolveColors(activeTheme, cfg.Theme["time"], "time")
	timeCfg := segment.TimeConfig{Format: cfg.Time.Format, FG: timeFG, BG: timeBG}

	available := map[string]segment.Segment{
		"path":    segment.NewPath(pathCfg),
		"git":     segment.NewGit(gitCfg),
		"kube":    segment.NewKube(kubeCfg),
		"prompt":  segment.NewPrompt(promptCfg),
		"time":    segment.NewTime(timeCfg),
		"newline": segment.NewNewline(),
	}

	var segments []segment.Segment
	for _, name := range cfg.Segments {
		seg, ok := available[name]
		if !ok {
			fmt.Fprintf(os.Stderr, "highline: unknown segment %q\n", name)
			os.Exit(1)
		}
		segments = append(segments, seg)
	}

	var sh renderer.Shell
	switch cfg.Shell {
	case "zsh":
		sh = renderer.ShellZsh
	default:
		sh = renderer.ShellBash
	}

	r := renderer.New(sh)
	fmt.Print(r.Render(segments))
}

// resolveColors returns fg/bg pointers by applying the built-in theme palette
// first, then allowing per-segment config overrides to win.
func resolveColors(thm theme.Theme, override config.ColorTheme, segName string) (fg *int, bg *int) {
	if p, ok := thm[segName]; ok {
		fg = &p.FG
		bg = &p.BG
	}
	if override.FG != nil {
		fg = override.FG
	}
	if override.BG != nil {
		bg = override.BG
	}
	return
}
