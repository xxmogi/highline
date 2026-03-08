package segment

import (
	"os"

	"github.com/user/highline/internal/color"
)

// PromptConfig holds configuration for PromptSegment.
type PromptConfig struct {
	Shell  string     // "zsh" or "bash" (default)
	FG     *int       // nil = color.BrightWhite (15)
	BG     *int       // nil = color.Black (0)
	UIDFn  func() int // nil = os.Getuid; injectable for tests
}

// PromptSegment renders the shell prompt symbol ($, %, or #).
type PromptSegment struct {
	cfg PromptConfig
}

func NewPrompt(cfg PromptConfig) *PromptSegment { return &PromptSegment{cfg: cfg} }

func (p *PromptSegment) Name() string { return "prompt" }

func (p *PromptSegment) Render() (string, int, int, error) {
	uidFn := os.Getuid
	if p.cfg.UIDFn != nil {
		uidFn = p.cfg.UIDFn
	}

	var sym string
	if uidFn() == 0 {
		sym = "#"
	} else if p.cfg.Shell == "zsh" {
		sym = "%"
	} else {
		sym = "$"
	}

	fg := color.BrightWhite
	if p.cfg.FG != nil {
		fg = *p.cfg.FG
	}
	bg := color.Black
	if p.cfg.BG != nil {
		bg = *p.cfg.BG
	}

	return sym, fg, bg, nil
}
