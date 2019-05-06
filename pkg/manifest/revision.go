package manifest

// Revision is the step before applying the next version of a manifest and
// potentially deleting leftovers from the old version. A revision with nil
// Next is considered as a deletion of all resources defined in the manifest.
type Revision struct {
	Prev *Manifest
	Next *Manifest
}

// HasNext returns false if the next manifest is not present. This indicates
// that the manifest should be deleted from the cluster using its previous
// state.
func (r Revision) HasNext() bool {
	return r.Next != nil
}

// CreateRevisions takes two slices of manifests and pairs matching
// manifests into revisions with previous and next manifest.
func CreateRevisions(prev, next []*Manifest) []Revision {
	revisions := make([]Revision, 0)

	for _, p := range prev {
		r := Revision{Prev: p}

		if i := findMatchingManifest(next, p); i >= 0 {
			r.Next = next[i]
		}

		revisions = append(revisions, r)
	}

	for _, n := range next {
		if i := findMatchingManifest(prev, n); i >= 0 {
			continue
		}

		revisions = append(revisions, Revision{Next: n})
	}

	return revisions
}

func findMatchingManifest(haystack []*Manifest, needle *Manifest) int {
	for i, m := range haystack {
		if m.Matches(needle) {
			return i
		}
	}

	return -1
}
