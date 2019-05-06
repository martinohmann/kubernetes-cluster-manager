package kubernetes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResourceSelectorMatches(t *testing.T) {
	cases := []struct {
		name    string
		a       ResourceSelector
		b       ResourceSelector
		matches bool
	}{
		{
			name:    "empty selectors",
			matches: true,
		},
		{
			name: "kind mismatch",
			a:    ResourceSelector{Kind: "Pod"},
			b:    ResourceSelector{Kind: "Deployment"},
		},
		{
			name: "namespace mismatch",
			a:    ResourceSelector{Kind: "Pod", Namespace: "default"},
			b:    ResourceSelector{Kind: "Pod", Namespace: "kube-system"},
		},
		{
			name: "name mismatch",
			a:    ResourceSelector{Kind: "Pod", Namespace: "default", Name: "foo"},
			b:    ResourceSelector{Kind: "Pod", Namespace: "default", Name: "bar"},
		},
		{
			name: "name mismatch, one empty",
			a:    ResourceSelector{Kind: "Pod", Namespace: "default", Name: ""},
			b:    ResourceSelector{Kind: "Pod", Namespace: "default", Name: "bar"},
		},
		{
			name:    "name match",
			a:       ResourceSelector{Kind: "Pod", Namespace: "default", Name: "foo"},
			b:       ResourceSelector{Kind: "Pod", Namespace: "default", Name: "foo"},
			matches: true,
		},
		{
			name: "empty names, label mismatch",
			a: ResourceSelector{
				Kind:      "Pod",
				Namespace: "default",
				Labels:    map[string]string{"app.kubernetes.io/name": "foo"},
			},
			b: ResourceSelector{
				Kind:      "Pod",
				Namespace: "default",
				Labels:    map[string]string{"app.kubernetes.io/name": "bar"},
			},
			matches: false,
		},
		{
			name: "empty names, label match",
			a: ResourceSelector{
				Kind:      "Pod",
				Namespace: "default",
				Labels:    map[string]string{"app.kubernetes.io/name": "foo"},
			},
			b: ResourceSelector{
				Kind:      "Pod",
				Namespace: "default",
				Labels:    map[string]string{"app.kubernetes.io/name": "foo"},
			},
			matches: true,
		},
		{
			name: "names take precedence over labels",
			a: ResourceSelector{
				Kind:      "Pod",
				Namespace: "default",
				Name:      "foo",
				Labels:    map[string]string{"app.kubernetes.io/name": "foo"},
			},
			b: ResourceSelector{
				Kind:      "Pod",
				Namespace: "default",
				Name:      "bar",
				Labels:    map[string]string{"app.kubernetes.io/name": "foo"},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.matches, tc.a.Matches(tc.b))
		})
	}
}
