package hook

import "github.com/martinohmann/kubernetes-cluster-manager/pkg/resource"

// Slice is a slice if hooks.
type Slice []*Hook

// Sort sorts the hooks slice.
func (s Slice) Sort() Slice {
	return sortHooks(s)
}

// Resources returns the hook resources.
func (s Slice) Resources() resource.Slice {
	rs := make(resource.Slice, len(s))

	for i, h := range s {
		rs[i] = h.Resource
	}

	return rs
}
