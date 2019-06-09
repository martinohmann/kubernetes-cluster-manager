package diff

import "github.com/martinohmann/go-difflib/difflib"

// Options is a set of diff options.
type Options struct {
	Filename string
	A, B     []byte
}

// Diff creates a diff based on o.
func Diff(o Options) string {
	unifiedDiff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(string(o.A)),
		B:        difflib.SplitLines(string(o.B)),
		FromFile: o.Filename,
		ToFile:   o.Filename,
		Context:  5,
		Color:    true,
	}

	out, _ := difflib.GetUnifiedDiffString(unifiedDiff)

	return out
}
