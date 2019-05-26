package hook

import (
	"testing"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/resource"
	"github.com/stretchr/testify/assert"
)

var unsorted = SliceMap{
	TypePreCreate: Slice{
		{Policy: "baz", Resource: &resource.Resource{Kind: "Job", Name: "foo"}},
		{Policy: "foo", Resource: &resource.Resource{Kind: "Job", Name: "baz"}},
		{Policy: "bar", Resource: &resource.Resource{Kind: "Job", Name: "foo"}},
	},
}

func TestSlice_Sort(t *testing.T) {
	expected := SliceMap{
		TypePreCreate: Slice{
			{Policy: "foo", Resource: &resource.Resource{Kind: "Job", Name: "baz"}},
			{Policy: "bar", Resource: &resource.Resource{Kind: "Job", Name: "foo"}},
			{Policy: "baz", Resource: &resource.Resource{Kind: "Job", Name: "foo"}},
		},
	}

	assert.Equal(t, expected, unsorted.SortSlices())
}
