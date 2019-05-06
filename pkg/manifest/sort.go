package manifest

// ByName implements sort.Interface and sorts a slice of *Manifest by Name.
type ByName []*Manifest

// Len implements Len from sort.Interface.
func (n ByName) Len() int {
	return len(n)
}

// Swap implements Swap from sort.Interface.
func (n ByName) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}

// Less implements Less from sort.Interface.
func (n ByName) Less(i, j int) bool {
	if n[i] == nil {
		return true
	}

	if n[j] == nil {
		return false
	}

	return n[i].Name < n[j].Name
}
