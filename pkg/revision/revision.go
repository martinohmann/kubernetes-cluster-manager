package revision

import (
	"bytes"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/hook"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/manifest"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/resource"
)

// Revision is the step before applying the next version of a manifest and
// potentially deleting leftovers from the old version. A revision with nil
// Next is considered as a deletion of all resources defined in the manifest.
// Consequently, a revision with a nil Current is considered an initial
// resource creation and will apply all resources.
type Revision struct {
	Current *manifest.Manifest
	Next    *manifest.Manifest
}

// ChangeSet is a container for resources that are sorted into buckets. These
// buckets help in finding the best upgrade strategy for a given manifest.
type ChangeSet struct {
	Revision           *Revision
	AddedResources     resource.Slice
	UpdatedResources   resource.Slice
	UnchangedResources resource.Slice
	RemovedResources   resource.Slice

	Hooks hook.SliceMap
}

// HasResourceChanges returns true if there are resource changes waiting to be
// applied. Resource changes are resource additions, deletions and updates.
func (c *ChangeSet) HasResourceChanges() bool {
	return len(c.AddedResources) > 0 || len(c.UpdatedResources) > 0 || len(c.RemovedResources) > 0
}

// Slice is a slice of revisions.
type Slice []*Revision

// Reverse reverses the order of a slice of *Revision. This is necessary to
// allow iterating all revisions in reverse order while deleting all manifests.
func (s Slice) Reverse() Slice {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}

	return s
}

// IsInitial returns true if r contains a new manifest, meaning that there is
// no current revision.
func (r *Revision) IsInitial() bool {
	return r.Current == nil && r.Next != nil
}

// IsRemoval returns true if r does not have a next manifest. This denotes that
// the manifest should be deleted from the cluster using the current revision.
func (r *Revision) IsRemoval() bool {
	return r.Current != nil && r.Next == nil
}

// IsUpgrade returns true if the manifest still exists in the next revision.
func (r *Revision) IsUpgrade() bool {
	return r.Current != nil && r.Next != nil
}

// IsValid returns true if there is at least a current or a next manifest in
// the revision.
func (r *Revision) IsValid() bool {
	return r.Current != nil || r.Next != nil
}

// Manifest returns the most recent manifest in the revision. If Next is
// present it will be returned, Current otherwise.
func (r *Revision) Manifest() *manifest.Manifest {
	if r.Next != nil {
		return r.Next
	}

	return r.Current
}

// NewSlice takes two slices of manifests and pairs matching
// manifests into revisions with current and next manifest.
func NewSlice(current, next []*manifest.Manifest) Slice {
	revisions := make(Slice, 0)

	for _, c := range current {
		r := &Revision{Current: c}

		if n, ok := manifest.FindMatching(next, c); ok {
			r.Next = n
		}

		revisions = append(revisions, r)
	}

	for _, n := range next {
		if _, ok := manifest.FindMatching(current, n); !ok {
			revisions = append(revisions, &Revision{Next: n})
		}
	}

	return revisions
}

// ChangeSet creates a ChangeSet for r. The change set categorizes resources
// into buckets (e.g. added, updated, unchanged, removed) and also contains the
// most recent hooks for this revision.
func (r *Revision) ChangeSet() *ChangeSet {
	if r.IsRemoval() {
		return &ChangeSet{
			Revision:         r,
			RemovedResources: r.Current.Resources,
			Hooks:            r.Current.Hooks,
		}
	}

	if r.IsInitial() {
		return &ChangeSet{
			Revision:       r,
			AddedResources: r.Next.Resources,
			Hooks:          r.Next.Hooks,
		}
	}

	c := &ChangeSet{
		Revision: r,
		Hooks:    r.Next.Hooks,
	}

	for _, current := range r.Current.Resources {
		res, ok := resource.FindMatching(r.Next.Resources, current)
		if !ok {
			c.RemovedResources = append(c.RemovedResources, current)
		} else if bytes.Compare(current.Content, res.Content) == 0 {
			c.UnchangedResources = append(c.UnchangedResources, res)
		} else {
			c.UpdatedResources = append(c.UpdatedResources, res)
		}
	}

	for _, next := range r.Next.Resources {
		_, ok := resource.FindMatching(r.Current.Resources, next)
		if !ok {
			c.AddedResources = append(c.AddedResources, next)
		}
	}

	return c
}
