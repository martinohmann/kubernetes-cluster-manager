package hook

import (
	"bytes"
	"sort"
)

type SliceMap map[Type]Slice

func (m SliceMap) Get(typ Type) Slice {
	return m[typ]
}

func (m SliceMap) Has(typ Type) bool {
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
		t := Type(k)
		buf.Write(m[t].Resources().Bytes())
	}

	return buf.Bytes()
}

func (m SliceMap) Sort() SliceMap {
	for _, v := range m {
		v.Sort()
	}

	return m
}
