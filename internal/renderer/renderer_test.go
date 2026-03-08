package renderer

import (
	"errors"
	"strings"
	"testing"

	"github.com/user/highline/internal/segment"
)

// stub implements segment.Segment for testing.
type stub struct {
	name string
	text string
	fg   int
	bg   int
	err  error
}

func (s *stub) Name() string                         { return s.name }
func (s *stub) Render() (string, int, int, error) { return s.text, s.fg, s.bg, s.err }

func segs(ss ...*stub) []segment.Segment {
	out := make([]segment.Segment, len(ss))
	for i, s := range ss {
		out[i] = s
	}
	return out
}

func TestRender_Empty(t *testing.T) {
	r := New(ShellBash)
	if got := r.Render(nil); got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestRender_SkipsErrorSegment(t *testing.T) {
	r := New(ShellBash)
	got := r.Render(segs(
		&stub{name: "ok", text: "hello", fg: ColorWhite, bg: ColorBlue},
		&stub{name: "bad", err: errors.New("no git")},
	))
	if !strings.Contains(got, "hello") {
		t.Errorf("expected 'hello' in output, got %q", got)
	}
}

func TestRender_AllErrorSegments(t *testing.T) {
	r := New(ShellBash)
	got := r.Render(segs(
		&stub{name: "bad", err: errors.New("skip me")},
	))
	if got != "" {
		t.Errorf("expected empty string when all segments error, got %q", got)
	}
}

func TestRender_BashEscapes(t *testing.T) {
	r := New(ShellBash)
	got := r.Render(segs(&stub{name: "p", text: "~/work", fg: ColorWhite, bg: ColorBlue}))
	if !strings.Contains(got, "~/work") {
		t.Errorf("expected text in output, got %q", got)
	}
	if !strings.Contains(got, "\001") || !strings.Contains(got, "\002") {
		t.Errorf("expected bash non-printing markers, got %q", got)
	}
	if strings.Contains(got, "%{") {
		t.Errorf("unexpected zsh markers in bash output")
	}
}

func TestRender_ZshEscapes(t *testing.T) {
	r := New(ShellZsh)
	got := r.Render(segs(&stub{name: "p", text: "~/work", fg: ColorWhite, bg: ColorBlue}))
	if !strings.Contains(got, "~/work") {
		t.Errorf("expected text in output, got %q", got)
	}
	if !strings.Contains(got, "%{") || !strings.Contains(got, "%}") {
		t.Errorf("expected zsh escape markers, got %q", got)
	}
	if strings.Contains(got, "\001") {
		t.Errorf("unexpected bash markers in zsh output")
	}
}

func TestRender_NewlineSegment(t *testing.T) {
	r := New(ShellBash)
	got := r.Render(segs(
		&stub{name: "a", text: "foo", fg: ColorWhite, bg: ColorBlue},
		&stub{name: "newline"},
		&stub{name: "b", text: "bar", fg: ColorBlack, bg: ColorGreen},
	))
	if !strings.Contains(got, "foo") || !strings.Contains(got, "bar") {
		t.Errorf("expected both segment texts, got %q", got)
	}
	if !strings.Contains(got, "\n") {
		t.Errorf("expected newline in output, got %q", got)
	}
	idx := strings.Index(got, "foo")
	idxNL := strings.Index(got, "\n")
	idxBar := strings.Index(got, "bar")
	if !(idx < idxNL && idxNL < idxBar) {
		t.Errorf("expected foo before newline before bar, got %q", got)
	}
}

func TestRender_MultipleSegmentsHaveSeparator(t *testing.T) {
	r := New(ShellBash)
	got := r.Render(segs(
		&stub{name: "a", text: "foo", fg: ColorWhite, bg: ColorBlue},
		&stub{name: "b", text: "bar", fg: ColorBlack, bg: ColorGreen},
	))
	if !strings.Contains(got, "foo") || !strings.Contains(got, "bar") {
		t.Errorf("expected both segment texts in output, got %q", got)
	}
	if !strings.Contains(got, separator) {
		t.Errorf("expected separator glyph in output, got %q", got)
	}
}
