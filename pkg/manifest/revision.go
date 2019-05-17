package manifest

import "bytes"

// Revision is the step before applying the next version of a manifest and
// potentially deleting leftovers from the old version. A revision with nil
// Next is considered as a deletion of all resources defined in the manifest.
type Revision struct {
	Current *Manifest
	Next    *Manifest
}

type ChangeSet struct {
	Revision           *Revision
	AddedResources     ResourceSlice
	ChangedResources   ResourceSlice
	UnchangedResources ResourceSlice
	RemovedResources   ResourceSlice

	Hooks HookSliceMap
}

type RevisionSlice []*Revision

// Reverse reverses the order of a slice of *Revision. This is necessary to
// allow iterating all revisions in reverse order while deleting all manifests.
func (s RevisionSlice) Reverse() RevisionSlice {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}

	return s
}

func (r *Revision) IsInstall() bool {
	return r.Current == nil
}

func (r *Revision) IsDelete() bool {
	return r.Next == nil
}

func (r *Revision) IsUpgrade() bool {
	return !r.IsInstall() && !r.IsDelete()
}

// CreateRevisions takes two slices of manifests and pairs matching
// manifests into revisions with previous and next manifest.
func CreateRevisions(current, next []*Manifest) RevisionSlice {
	revisions := make(RevisionSlice, 0)

	for _, c := range current {
		r := &Revision{Current: c}

		if n, ok := findMatchingManifest(next, c); ok {
			r.Next = n
		}

		revisions = append(revisions, r)
	}

	for _, n := range next {
		if _, ok := findMatchingManifest(current, n); !ok {
			revisions = append(revisions, &Revision{Next: n})
		}
	}

	return revisions
}

// ChangeSet creates the change set for r. The change set categorizes resources
// into buckets (e.g. added, changed, unchanged, removed) and also contains the
// most recent hooks for this revision.
func (r *Revision) ChangeSet() *ChangeSet {
	if r.IsDelete() {
		return &ChangeSet{
			Revision:         r,
			RemovedResources: r.Current.resources,
			Hooks:            r.Current.hooks,
		}
	}

	if r.IsInstall() {
		return &ChangeSet{
			Revision:       r,
			AddedResources: r.Next.resources,
			Hooks:          r.Next.hooks,
		}

	}

	c := &ChangeSet{
		Revision:           r,
		AddedResources:     make(ResourceSlice, 0),
		ChangedResources:   make(ResourceSlice, 0),
		UnchangedResources: make(ResourceSlice, 0),
		RemovedResources:   make(ResourceSlice, 0),
		Hooks:              r.Next.hooks,
	}

	for _, current := range r.Current.resources {
		res, ok := findMatchingResource(r.Next.resources, current)
		if !ok {
			c.RemovedResources = append(c.RemovedResources, current)
		} else if bytes.Compare(current.Content, res.Content) == 0 {
			c.UnchangedResources = append(c.UnchangedResources, res)
		} else {
			c.ChangedResources = append(c.ChangedResources, res)
		}
	}

	for _, next := range r.Next.resources {
		_, ok := findMatchingResource(r.Current.resources, next)
		if !ok {
			c.AddedResources = append(c.AddedResources, next)
		}
	}

	return c
}

func findMatchingManifest(haystack []*Manifest, needle *Manifest) (*Manifest, bool) {
	for _, m := range haystack {
		if m.Name == needle.Name {
			return m, true
		}
	}

	return nil, false
}

func findMatchingResource(haystack []*Resource, needle *Resource) (*Resource, bool) {
	for _, r := range haystack {
		if r.matches(needle) {
			return r, true
		}
	}

	return nil, false
}