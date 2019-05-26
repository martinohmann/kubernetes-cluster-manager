package hook

import (
	"bytes"
	"sort"
)

// SliceMap is a map of hook slices.
type SliceMap map[string]Slice

// Bytes returns the raw resource bytes for all hooks in the map.
func (m SliceMap) Bytes() []byte {
	var buf bytes.Buffer

	keys := make([]string, 0, len(m))

	for k := range m {
		keys = append(keys, string(k))
	}

	sort.Strings(keys)

	for _, k := range keys {
		buf.Write(m[k].Resources().Bytes())
	}

	return buf.Bytes()
}

// SortSlices sorts all slices of the map.
func (m SliceMap) SortSlices() SliceMap {
	for _, v := range m {
		v.Sort()
	}

	return m
}
