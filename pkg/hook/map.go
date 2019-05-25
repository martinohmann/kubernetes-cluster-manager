package hook

import (
	"bytes"
	"sort"
)

type SliceMap map[string]Slice

func (m SliceMap) Get(typ string) Slice {
	return m[typ]
}

func (m SliceMap) Has(typ string) bool {
	hooks, ok := m[typ]

	return ok && len(hooks) > 0
}

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

func (m SliceMap) Sort() SliceMap {
	for _, v := range m {
		v.Sort()
	}

	return m
}
