package hook

import (
	"testing"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/resource"
	"github.com/stretchr/testify/assert"
)

var unsorted = SliceMap{
	TypePreCreate: Slice{
		{WaitFor: "condition=baz", Resource: &resource.Resource{Kind: "Job", Name: "foo"}},
		{WaitFor: "condition=foo", Resource: &resource.Resource{Kind: "Job", Name: "baz"}},
		{WaitFor: "condition=bar", Resource: &resource.Resource{Kind: "Job", Name: "foo"}},
	},
}

func TestSlice_Sort(t *testing.T) {
	expected := SliceMap{
		TypePreCreate: Slice{
			{WaitFor: "condition=foo", Resource: &resource.Resource{Kind: "Job", Name: "baz"}},
			{WaitFor: "condition=bar", Resource: &resource.Resource{Kind: "Job", Name: "foo"}},
			{WaitFor: "condition=baz", Resource: &resource.Resource{Kind: "Job", Name: "foo"}},
		},
	}

	assert.Equal(t, expected, unsorted.SortSlices())
}
