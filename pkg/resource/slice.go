package resource

import "strings"

// Slice is a slice of resources.
type Slice []*Resource

// Bytes converts the resource slice to raw bytes.
func (s Slice) Bytes() []byte {
	var buf Buffer

	for _, r := range s {
		buf.Write(r.Content)
	}

	return buf.Bytes()
}

// Sort sorts the slice in the given order.
func (s Slice) Sort(order ResourceOrder) Slice {
	return sortResources(s, order)
}

// String implements fmt.Stringer
func (s Slice) String() string {
	names := make([]string, len(s))

	for i, r := range s {
		names[i] = r.String()
	}

	return strings.Join(names, "\n")
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
