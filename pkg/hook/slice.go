package hook

import "github.com/martinohmann/kubernetes-cluster-manager/pkg/resource"

type Slice []*Hook

func (s Slice) Sort() Slice {
	return sortHooks(s)
}

func (s Slice) Resources() resource.Slice {
	rs := make(resource.Slice, len(s))

	for i, h := range s {
		rs[i] = h.Resource
	}

	return rs
}
