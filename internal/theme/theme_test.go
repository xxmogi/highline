package theme

import (
	"testing"
)

func TestGet_KnownTheme(t *testing.T) {
	tests := []struct {
		name    string
		segment string
		wantFG  int
		wantBG  int
	}{
		{"monokai", "path", 0, 10},
		{"monokai", "git", 0, 11},
		{"monokai", "kube", 0, 5},
		{"solarized-dark", "path", 15, 4},
		{"solarized-light", "path", 8, 3},
		{"nord", "git", 15, 12},
		{"dracula", "path", 15, 5},
	}
	for _, tt := range tests {
		t.Run(tt.name+"/"+tt.segment, func(t *testing.T) {
			thm, ok := Get(tt.name)
			if !ok {
				t.Fatalf("Get(%q) not found", tt.name)
			}
			p, ok := thm[tt.segment]
			if !ok {
				t.Fatalf("theme %q has no segment %q", tt.name, tt.segment)
			}
			if p.FG != tt.wantFG {
				t.Errorf("FG: want %d, got %d", tt.wantFG, p.FG)
			}
			if p.BG != tt.wantBG {
				t.Errorf("BG: want %d, got %d", tt.wantBG, p.BG)
			}
		})
	}
}

func TestGet_CaseInsensitive(t *testing.T) {
	_, ok := Get("Monokai")
	if !ok {
		t.Error("Get(\"Monokai\") should succeed (case-insensitive)")
	}
	_, ok = Get("NORD")
	if !ok {
		t.Error("Get(\"NORD\") should succeed (case-insensitive)")
	}
}

func TestGet_Unknown(t *testing.T) {
	_, ok := Get("nonexistent")
	if ok {
		t.Error("Get(\"nonexistent\") should return false")
	}
}

func TestNames(t *testing.T) {
	names := Names()
	if len(names) != 5 {
		t.Errorf("expected 5 theme names, got %d: %v", len(names), names)
	}
	// Verify sorted order.
	for i := 1; i < len(names); i++ {
		if names[i] < names[i-1] {
			t.Errorf("Names() not sorted: %v", names)
		}
	}
	// Verify all expected names present.
	expected := map[string]bool{
		"monokai": true, "solarized-dark": true, "solarized-light": true,
		"nord": true, "dracula": true,
	}
	for _, n := range names {
		if !expected[n] {
			t.Errorf("unexpected theme name: %q", n)
		}
	}
}
