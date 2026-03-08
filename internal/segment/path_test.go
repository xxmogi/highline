package segment

import (
	"testing"
)

func TestShortenPath(t *testing.T) {
	tests := []struct {
		input    string
		maxDepth int
		want     string
	}{
		{"~", 4, "~"},
		{"~/a", 4, "~/a"},
		{"~/a/b", 4, "~/a/b"},
		{"~/a/b/c", 4, "~/a/b/c"},
		{"~/a/b/c/d", 4, "~/.../d"},
		{"~/a/b/c/d/e", 4, "~/.../e"},
		{"/usr/local/bin", 4, "/usr/local/bin"},
		{"/a/b/c/d", 4, "/a/.../d"},
		// Custom maxDepth
		{"~/a/b/c", 3, "~/.../c"},
		{"~/a/b", 3, "~/a/b"},
	}

	for _, tc := range tests {
		got := shortenPath(tc.input, tc.maxDepth)
		if got != tc.want {
			t.Errorf("shortenPath(%q, %d) = %q, want %q", tc.input, tc.maxDepth, got, tc.want)
		}
	}
}

func TestPathSegment_Name(t *testing.T) {
	p := NewPath(PathConfig{})
	if p.Name() != "path" {
		t.Errorf("expected 'path', got %q", p.Name())
	}
}

func TestPathSegment_Render(t *testing.T) {
	p := NewPath(PathConfig{})
	text, fg, bg, err := p.Render()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if text == "" {
		t.Error("expected non-empty path")
	}
	if fg == 0 && bg == 0 {
		t.Error("expected non-zero color values")
	}
	_ = fg
	_ = bg
}

func TestPathSegment_Render_CustomColors(t *testing.T) {
	fg := 3
	bg := 5
	p := NewPath(PathConfig{FG: &fg, BG: &bg})
	_, gotFG, gotBG, err := p.Render()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotFG != fg {
		t.Errorf("expected fg=%d, got %d", fg, gotFG)
	}
	if gotBG != bg {
		t.Errorf("expected bg=%d, got %d", bg, gotBG)
	}
}

func TestPathSegment_Render_CustomMaxDepth(t *testing.T) {
	// MaxDepth=0 should default to 4 (no panic).
	p := NewPath(PathConfig{MaxDepth: 0})
	_, _, _, err := p.Render()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
