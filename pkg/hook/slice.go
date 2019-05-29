package hook

import (
	"strings"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/resource"
)

// Slice is a slice if hooks.
type Slice []*Hook

// Sort sorts the hooks slice.
func (s Slice) Sort() Slice {
	return sortHooks(s)
}

// String implements fmt.Stringer
func (s Slice) String() string {
	names := make([]string, len(s))

	for i, h := range s {
		names[i] = h.String()
	}

	return strings.Join(names, "\n")
}

// Resources returns the hook resources.
func (s Slice) Resources() resource.Slice {
	rs := make(resource.Slice, len(s))

	for i, h := range s {
		rs[i] = h.Resource
	}

	return rs
}
