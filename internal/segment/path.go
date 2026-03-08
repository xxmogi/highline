package segment

import (
	"os"
	"strings"

	"github.com/user/highline/internal/color"
)

// PathConfig holds optional color and depth configuration for PathSegment.
type PathConfig struct {
	FG       *int // nil = color.White
	BG       *int // nil = color.Blue
	MaxDepth int  // 0 = 4
}

// PathSegment renders the current working directory.
type PathSegment struct {
	cfg PathConfig
}

func NewPath(cfg PathConfig) *PathSegment { return &PathSegment{cfg: cfg} }

func (p *PathSegment) Name() string { return "path" }

func (p *PathSegment) Render() (string, int, int, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", 0, 0, err
	}

	home, err := os.UserHomeDir()
	if err == nil && strings.HasPrefix(cwd, home) {
		cwd = "~" + cwd[len(home):]
	}

	maxDepth := p.cfg.MaxDepth
	if maxDepth <= 0 {
		maxDepth = 4
	}
	cwd = shortenPath(cwd, maxDepth)

	fg := color.White
	if p.cfg.FG != nil {
		fg = *p.cfg.FG
	}
	bg := color.Blue
	if p.cfg.BG != nil {
		bg = *p.cfg.BG
	}

	return cwd, fg, bg, nil
}

// shortenPath abbreviates paths with more than maxDepth components.
// e.g. (maxDepth=4) ~/a/b/c/d → ~/.../d, /a/b/c/d → /a/.../d
func shortenPath(p string, maxDepth int) string {
	parts := strings.Split(p, "/")
	if len(parts) <= maxDepth {
		return p
	}
	// Absolute path: parts[0] is empty, parts[1] is the first real component.
	if parts[0] == "" {
		return "/" + parts[1] + "/.../" + parts[len(parts)-1]
	}
	return parts[0] + "/.../" + parts[len(parts)-1]
}
