package segment

// Segment represents a single unit of the PS1 prompt.
// If Render returns a non-nil error, the segment is skipped.
type Segment interface {
	Name() string
	Render() (text string, fg int, bg int, err error)
}
