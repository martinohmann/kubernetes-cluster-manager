package manifest

type ByName []*Manifest

// Len implements Len from sort.Interface.
func (m ByName) Len() int {
	return len(m)
}

// Swap implements Swap from sort.Interface.
func (m ByName) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

// Less implements Less from sort.Interface.
func (m ByName) Less(i, j int) bool {
	a, b := m[i], m[j]

	if a == nil {
		return true
	}

	if b == nil {
		return false
	}

	return a.Name < b.Name
}
