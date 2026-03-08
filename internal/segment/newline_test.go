package segment

import "testing"

func TestNewline_Name(t *testing.T) {
	n := NewNewline()
	if n.Name() != "newline" {
		t.Errorf("Name: want %q, got %q", "newline", n.Name())
	}
}

func TestNewline_RenderReturnsError(t *testing.T) {
	n := NewNewline()
	_, _, _, err := n.Render()
	if err == nil {
		t.Error("expected error from NewlineSegment.Render(), got nil")
	}
}
