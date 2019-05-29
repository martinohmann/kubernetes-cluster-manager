package hook

import "sort"

type hookSorter struct {
	hooks []*Hook
}

func newHookSorter(hooks []*Hook) *hookSorter {
	return &hookSorter{
		hooks: hooks,
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

	if a.Resource.Name == b.Resource.Name {
		return a.WaitFor < b.WaitFor
	}

	return a.Resource.Name < b.Resource.Name
}

func sortHooks(hooks []*Hook) []*Hook {
	s := newHookSorter(hooks)

	sort.Sort(s)

	return hooks
}
