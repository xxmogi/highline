package segment

import "errors"

// NewlineSegment is a marker segment that signals the renderer to insert a
// line break. Its Render method is never called directly by the renderer.
type NewlineSegment struct{}

func NewNewline() *NewlineSegment { return &NewlineSegment{} }

func (n *NewlineSegment) Name() string { return "newline" }

func (n *NewlineSegment) Render() (string, int, int, error) {
	return "", 0, 0, errors.New("newline segment: Render must not be called directly")
}
