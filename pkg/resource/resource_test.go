package resource

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResource_matches(t *testing.T) {
	cases := []struct {
		description string
		a, b        *Resource
		matches     bool
	}{
		{
			description: "both nil",
			matches:     true,
		},
		{
			description: "a nil",
			b:           &Resource{Name: "foo"},
		},
		{
			description: "b nil",
			a:           &Resource{Name: "foo"},
		},
		{
			description: "different kind",
			a:           &Resource{Name: "foo", Kind: "Deployment"},
			b:           &Resource{Name: "foo", Kind: "Pod"},
		},
		{
			description: "different namespace",
			a:           &Resource{Name: "foo", Kind: "Pod", Namespace: "default"},
			b:           &Resource{Name: "foo", Kind: "Pod", Namespace: "kube-system"},
		},
		{
			description: "same name, kind and namespace",
			a:           &Resource{Name: "foo", Kind: "Pod", Namespace: "kube-system"},
			b:           &Resource{Name: "foo", Kind: "Pod", Namespace: "kube-system"},
			matches:     true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			assert.Equal(t, tc.matches, tc.a.matches(tc.b))
		})
	}
}
