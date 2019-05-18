// Adapted from https://github.com/helm/helm/blob/master/pkg/tiller/kind_sorter.go

package manifest

import "sort"

type hookSorter struct {
	order map[string]int
	hooks []*Hook
}

func newHookSorter(hooks []*Hook, order ResourceOrder) *hookSorter {
	o := make(map[string]int)

	for k, v := range order {
		o[v] = k
	}

	return &hookSorter{
		hooks: hooks,
		order: o,
	}
}

// Len implements Len from sort.Interface.
func (s *hookSorter) Len() int {
	return len(s.hooks)
}

// Swap implements Swap from sort.Interface.
func (s *hookSorter) Swap(i, j int) {
	s.hooks[i], s.hooks[j] = s.hooks[j], s.hooks[i]
}

// Less implements Less from sort.Interface.
func (s *hookSorter) Less(i, j int) bool {
	a, b := s.hooks[i], s.hooks[j]

	aPos, aok := s.order[a.Kind]
	bPos, bok := s.order[b.Kind]

	if !aok && !bok {
		if a.Kind == b.Kind {
			return a.Name < b.Name
		}

		return a.Kind < b.Kind
	}

	if !aok || !bok {
		return aok
	}

	if aPos == bPos {
		return a.Name < b.Name
	}

	return aPos < bPos
}

func sortHooks(hooks []*Hook, order ResourceOrder) []*Hook {
	s := newHookSorter(hooks, order)

	sort.Sort(s)

	return hooks
}
