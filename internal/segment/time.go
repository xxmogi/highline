package segment

import (
	"time"

	"github.com/user/highline/internal/color"
)

// TimeConfig holds configuration for TimeSegment.
type TimeConfig struct {
	Format string         // Go time.Format string; "" = "15:04"
	FG     *int           // nil = color.Black
	BG     *int           // nil = color.White
	NowFn  func() time.Time // nil = time.Now; injectable for tests
}

// TimeSegment renders the current time.
type TimeSegment struct {
	cfg TimeConfig
}

func NewTime(cfg TimeConfig) *TimeSegment { return &TimeSegment{cfg: cfg} }

func (t *TimeSegment) Name() string { return "time" }

func (t *TimeSegment) Render() (string, int, int, error) {
	now := time.Now()
	if t.cfg.NowFn != nil {
		now = t.cfg.NowFn()
	}

	format := "15:04"
	if t.cfg.Format != "" {
		format = t.cfg.Format
	}

	fg := color.Black
	if t.cfg.FG != nil {
		fg = *t.cfg.FG
	}
	bg := color.White
	if t.cfg.BG != nil {
		bg = *t.cfg.BG
	}

	return now.Format(format), fg, bg, nil
}
