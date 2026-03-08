package segment

import (
	"testing"
	"time"
)

func fixedTime(s string) func() time.Time {
	t, err := time.Parse("2006-01-02 15:04:05", s)
	if err != nil {
		panic(err)
	}
	return func() time.Time { return t }
}

func TestTime_DefaultFormat(t *testing.T) {
	seg := NewTime(TimeConfig{NowFn: fixedTime("2024-03-08 14:30:00")})
	text, _, _, err := seg.Render()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if text != "14:30" {
		t.Errorf("want %q, got %q", "14:30", text)
	}
}

func TestTime_CustomFormat(t *testing.T) {
	seg := NewTime(TimeConfig{
		Format: "2006-01-02",
		NowFn:  fixedTime("2024-03-08 14:30:00"),
	})
	text, _, _, err := seg.Render()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if text != "2024-03-08" {
		t.Errorf("want %q, got %q", "2024-03-08", text)
	}
}

func TestTime_DefaultColors(t *testing.T) {
	seg := NewTime(TimeConfig{NowFn: fixedTime("2024-03-08 09:00:00")})
	_, fg, bg, err := seg.Render()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fg != 0 {
		t.Errorf("default FG: want 0 (Black), got %d", fg)
	}
	if bg != 7 {
		t.Errorf("default BG: want 7 (White), got %d", bg)
	}
}

func TestTime_CustomColors(t *testing.T) {
	fg, bg := 3, 5
	seg := NewTime(TimeConfig{NowFn: fixedTime("2024-03-08 09:00:00"), FG: &fg, BG: &bg})
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

func TestTime_Name(t *testing.T) {
	seg := NewTime(TimeConfig{})
	if seg.Name() != "time" {
		t.Errorf("Name: want %q, got %q", "time", seg.Name())
	}
}
