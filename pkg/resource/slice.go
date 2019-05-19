package resource

import "github.com/martinohmann/kubernetes-cluster-manager/pkg/kubernetes"

type Slice []*Resource

// Bytes converts the resource slice to raw bytes.
func (s Slice) Bytes() []byte {
	var buf buffer

	for _, r := range s {
		buf.Write(r.Content)
	}

	return buf.Bytes()
}

// Selectors creates a kubernetes.ResourceSelector for each Resource in s.
func (s Slice) Selectors() []kubernetes.ResourceSelector {
	rs := make([]kubernetes.ResourceSelector, 0)

	for _, res := range s {
		rs = append(rs, res.Selector())
	}

	return rs
}

// Sort sorts the slice in the given order.
func (s Slice) Sort(order ResourceOrder) Slice {
	return sortResources(s, order)
}

// FindMatching searches haystack for a resource matching needle and returns it
// if found, nil otherwise.
func FindMatching(haystack []*Resource, needle *Resource) (*Resource, bool) {
	for _, r := range haystack {
		if r.matches(needle) {
			return r, true
		}
	}

	return nil, false
}
