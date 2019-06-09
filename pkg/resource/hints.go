package resource

// Hint is the type for resource hints. Hints are used to add context to a
// resource so that output formatters can make a decision on how to format it.
type Hint int

const (
	// NoChange indicates that the resource has not changed.
	NoChange Hint = iota
	// Addition indicates that the resource was newly added.
	Addition
	// Update indicates that the resource content changed.
	Update
	// Removal indicates that the resource will be removed.
	Removal
)

// String implements fmt.Stringer
func (h Hint) String() string {
	switch h {
	case NoChange:
		return "no change"
	case Addition:
		return "addition"
	case Update:
		return "update"
	case Removal:
		return "removal"
	}

	return "unknown"
}

// WithHint sets a hint on the resource.
func (r *Resource) WithHint(hint Hint) *Resource {
	r.hint = hint

	return r
}

// WithContentHint hints the resource with its current content so that
// diffs can be generated for the new content in formatted output.
func (r *Resource) WithContentHint(content []byte) *Resource {
	r.contentHint = content

	return r
}

// WithHint sets a hint on all resources in the slice.
func (s Slice) WithHint(hint Hint) Slice {
	for _, r := range s {
		r.WithHint(hint)
	}

	return s
}
