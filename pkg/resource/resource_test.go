package resource

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindMatching(t *testing.T) {
	cases := []struct {
		description string
		r           *Resource
		expected    bool
	}{
		{
			description: "both nil",
			expected:    true,
		},
		{
			description: "a nil",
			r:           &Resource{Name: "foo"},
		},
		{
			description: "b nil",
			r:           &Resource{Name: "foo"},
		},
		{
			description: "different kind",
			r:           &Resource{Name: "foo", Kind: "Deployment"},
		},
		{
			description: "different namespace",
			r:           &Resource{Name: "foo", Kind: "Pod", Namespace: "default"},
		},
		{
			description: "same name, kind and namespace",
			r:           &Resource{Name: "foo", Kind: "Pod", Namespace: "kube-system"},
			expected:    true,
		},
	}

	resources := []*Resource{
		nil,
		{Name: "foo", Kind: "Pod", Namespace: "kube-system"},
		New(nil, Head{Kind: "Pod", Metadata: Metadata{Name: "foo", Namespace: "kube-system"}}),
	}

	_ = resources

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			r, ok := FindMatching(resources, tc.r)

			assert.Equal(t, tc.expected, ok)

			if tc.expected {
				assert.Equal(t, tc.r, r)
			}
		})
	}
}
