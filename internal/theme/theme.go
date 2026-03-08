package theme

import (
	"sort"
	"strings"
)

// Palette is a fg/bg color pair for one segment.
type Palette struct {
	FG int
	BG int
}

// Theme maps segment name → Palette.
type Theme map[string]Palette

var builtin = map[string]Theme{
	"monokai": {
		"path": {FG: 0, BG: 10},
		"git":  {FG: 0, BG: 11},
		"kube": {FG: 0, BG: 5},
	},
	"solarized-dark": {
		"path": {FG: 15, BG: 4},
		"git":  {FG: 8, BG: 2},
		"kube": {FG: 8, BG: 6},
	},
	"solarized-light": {
		"path": {FG: 8, BG: 3},
		"git":  {FG: 8, BG: 2},
		"kube": {FG: 8, BG: 6},
	},
	"nord": {
		"path": {FG: 15, BG: 4},
		"git":  {FG: 15, BG: 12},
		"kube": {FG: 8, BG: 6},
	},
	"dracula": {
		"path": {FG: 15, BG: 5},
		"git":  {FG: 0, BG: 10},
		"kube": {FG: 15, BG: 12},
	},
}

// Get returns the built-in theme for the given name (case-insensitive).
func Get(name string) (Theme, bool) {
	t, ok := builtin[strings.ToLower(name)]
	return t, ok
}

// Names returns all built-in theme names in sorted order.
func Names() []string {
	names := make([]string, 0, len(builtin))
	for k := range builtin {
		names = append(names, k)
	}
	sort.Strings(names)
	return names
}
