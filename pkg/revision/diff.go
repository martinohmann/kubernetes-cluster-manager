package revision

import "github.com/martinohmann/go-difflib/difflib"

// Diff creates a git style diff for the revision.
func (r *Revision) Diff() string {
	var a, b, fromFile, toFile string

	if r.Current != nil {
		fromFile = r.Current.Filename()
		a = string(r.Current.Content())
	}

	if r.Next != nil {
		toFile = r.Next.Filename()
		b = string(r.Next.Content())
	}

	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(a),
		B:        difflib.SplitLines(b),
		FromFile: fromFile,
		ToFile:   toFile,
		Context:  5,
		Color:    true,
	}

	out, _ := difflib.GetUnifiedDiffString(diff)

	return out
}
