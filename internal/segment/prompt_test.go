package segment

import "testing"

func uidFn(uid int) func() int { return func() int { return uid } }

func TestPrompt_Name(t *testing.T) {
	seg := NewPrompt(PromptConfig{})
	if seg.Name() != "prompt" {
		t.Errorf("Name: want %q, got %q", "prompt", seg.Name())
	}
}

func TestPrompt_Symbols(t *testing.T) {
	tests := []struct {
		name    string
		uid     int
		shell   string
		wantSym string
	}{
		{"bash non-root", 1000, "bash", "$"},
		{"bash root",     0,    "bash", "#"},
		{"zsh non-root",  1000, "zsh",  "%"},
		{"zsh root",      0,    "zsh",  "#"},
		{"default shell", 1000, "",     "$"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			seg := NewPrompt(PromptConfig{Shell: tt.shell, UIDFn: uidFn(tt.uid)})
			sym, _, _, err := seg.Render()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if sym != tt.wantSym {
				t.Errorf("want %q, got %q", tt.wantSym, sym)
			}
		})
	}
}

func TestPrompt_DefaultColors(t *testing.T) {
	seg := NewPrompt(PromptConfig{UIDFn: uidFn(1000)})
	_, fg, bg, err := seg.Render()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fg != 15 {
		t.Errorf("default FG: want 15 (BrightWhite), got %d", fg)
	}
	if bg != 0 {
		t.Errorf("default BG: want 0 (Black), got %d", bg)
	}
}

func TestPrompt_CustomColors(t *testing.T) {
	fg, bg := 2, 8
	seg := NewPrompt(PromptConfig{UIDFn: uidFn(1000), FG: &fg, BG: &bg})
	_, gotFG, gotBG, err := seg.Render()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotFG != fg {
		t.Errorf("FG: want %d, got %d", fg, gotFG)
	}
	if gotBG != bg {
		t.Errorf("BG: want %d, got %d", bg, gotBG)
	}
}
