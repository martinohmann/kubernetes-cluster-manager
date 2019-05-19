package resource

type Slice []*Resource

// Bytes converts the resource slice to raw bytes.
func (s Slice) Bytes() []byte {
	var buf buffer

	for _, r := range s {
		buf.Write(r.Content)
	}

	return buf.Bytes()
}

// Sort sorts the slice in the given order.
func (s Slice) Sort(order ResourceOrder) Slice {
	return sortResources(s, order)
}

// FindMatching searches haystack for a resource matching needle and returns it
// if found, nil otherwise.
func FindMatching(haystack []*Resource, needle *Resource) (*Resource, bool) {
	for _, r := range haystack {
		if r.matches(needle) {
			return r, true
		}
	}

	return nil, false
}
