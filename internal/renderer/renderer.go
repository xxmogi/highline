package renderer

import (
	"fmt"
	"strings"

	"github.com/user/highline/internal/segment"
)

// Shell represents the target shell for prompt escaping.
type Shell string

const (
	ShellBash Shell = "bash"
	ShellZsh  Shell = "zsh"
)

// separator is the Powerline right-arrow glyph (U+E0B0).
const separator = "\ue0b0"

// Renderer builds a Powerline-style prompt string.
type Renderer struct {
	shell Shell
}

// New returns a new Renderer for the given shell.
func New(shell Shell) *Renderer {
	return &Renderer{shell: shell}
}

type renderedPart struct {
	text string
	fg   int
	bg   int
}

// Render renders all segments into a prompt string.
// "newline" segments cause a line break; segments returning an error are skipped.
func (r *Renderer) Render(segments []segment.Segment) string {
	var sb strings.Builder
	var lineParts []renderedPart

	flush := func() {
		if len(lineParts) > 0 {
			sb.WriteString(r.renderLine(lineParts))
			lineParts = nil
		}
	}

	for _, seg := range segments {
		if seg.Name() == "newline" {
			flush()
			sb.WriteString("\n")
			continue
		}
		text, fg, bg, err := seg.Render()
		if err != nil {
			continue
		}
		lineParts = append(lineParts, renderedPart{text, fg, bg})
	}
	flush()
	return sb.String()
}

// renderLine renders one Powerline line from a slice of already-rendered parts.
func (r *Renderer) renderLine(parts []renderedPart) string {
	if len(parts) == 0 {
		return ""
	}

	var sb strings.Builder
	for i, p := range parts {
		// Set BG and FG for the segment content.
		sb.WriteString(r.escape(fmt.Sprintf("\033[38;5;%dm\033[48;5;%dm", p.fg, p.bg)))
		sb.WriteString(" ")
		sb.WriteString(p.text)
		sb.WriteString(" ")

		// Separator: FG = current BG, BG = next BG (or reset).
		if i+1 < len(parts) {
			nextBG := parts[i+1].bg
			sb.WriteString(r.escape(fmt.Sprintf("\033[38;5;%dm\033[48;5;%dm", p.bg, nextBG)))
		} else {
			sb.WriteString(r.escape(fmt.Sprintf("\033[38;5;%dm\033[49m", p.bg)))
		}
		sb.WriteString(separator)
	}

	// Reset all attributes.
	sb.WriteString(r.escape("\033[0m"))
	return sb.String()
}

// escape wraps an ANSI escape sequence in shell-specific non-printing markers.
func (r *Renderer) escape(ansi string) string {
	switch r.shell {
	case ShellZsh:
		return "%{" + ansi + "%}"
	default: // bash
		return "\001" + ansi + "\002"
	}
}
