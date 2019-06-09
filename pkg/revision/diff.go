package revision

import "github.com/martinohmann/kubernetes-cluster-manager/pkg/diff"

// DiffOptions returns the diff.Options for this revision.
func (r *Revision) DiffOptions() diff.Options {
	var o diff.Options

	if r.Current != nil {
		o.A = r.Current.Content()
		o.Filename = r.Current.Filename()
	}

	if r.Next != nil {
		o.B = r.Next.Content()
		o.Filename = r.Next.Filename()
	}

	return o
}
