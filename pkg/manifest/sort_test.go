package manifest

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSortByName(t *testing.T) {
	manifests := []*Manifest{
		{Name: "xyz"},
		nil,
		{Name: "abc"},
		nil,
	}

	expected := []*Manifest{
		nil,
		nil,
		{Name: "abc"},
		{Name: "xyz"},
	}

	sort.Sort(ByName(manifests))

	assert.Equal(t, expected, manifests)
}
