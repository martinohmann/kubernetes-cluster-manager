package manifest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestManifestResources(t *testing.T) {
	content := []byte(`---
foo: bar
---
kind: Pod
metadata:
  name: foo
  namespace: kube-system

---
kind: Deployment
metadata:
  name: bar
  namespace: apps

---
kind: ConfigMap
metadata:
  name: baz
data:
  SOME_ENV: somevalue

---
  
---`)

	expected := []ResourceSelector{
		{
			Kind:      "Pod",
			Name:      "foo",
			Namespace: "kube-system",
		},
		{
			Kind:      "Deployment",
			Name:      "bar",
			Namespace: "apps",
		},
		{
			Kind: "ConfigMap",
			Name: "baz",
		},
	}

	m := &Manifest{Name: "manifest", Content: content}

	actual := m.Resources()

	assert.Equal(t, expected, actual)
}

func TestRevisionGetVanishedResources(t *testing.T) {
	cases := []struct {
		name     string
		revision Revision
		expected []ResourceSelector
	}{
		{
			name: "nothing removed",
			revision: Revision{
				Prev: &Manifest{Content: []byte(`---
kind: Pod
metadata:
  name: foo
  namespace: bar
`)},
				Next: &Manifest{Content: []byte(`---
kind: Pod
metadata:
  name: foo
  namespace: bar
`)},
			},
			expected: []ResourceSelector{},
		},
		{
			name: "pod removed",
			revision: Revision{
				Prev: &Manifest{Content: []byte(`---
kind: Pod
metadata:
  name: foo
  namespace: bar
`)},
				Next: &Manifest{Content: []byte(`---
`)},
			},
			expected: []ResourceSelector{
				{
					Kind:      "Pod",
					Name:      "foo",
					Namespace: "bar",
				},
			},
		},
		{
			name: "pod moved to another namespace",
			revision: Revision{
				Prev: &Manifest{Content: []byte(`---
kind: Pod
metadata:
  name: foo
  namespace: bar
`)},
				Next: &Manifest{Content: []byte(`---
kind: Pod
metadata:
  name: foo
  namespace: baz
`)},
			},
			expected: []ResourceSelector{
				{
					Kind:      "Pod",
					Name:      "foo",
					Namespace: "bar",
				},
			},
		},
		{
			name: "pod renamed",
			revision: Revision{
				Prev: &Manifest{Content: []byte(`---
kind: Pod
metadata:
  name: foo
  namespace: bar
`)},
				Next: &Manifest{Content: []byte(`---
kind: Pod
metadata:
  name: foo2
  namespace: bar
`)},
			},
			expected: []ResourceSelector{
				{
					Kind:      "Pod",
					Name:      "foo",
					Namespace: "bar",
				},
			},
		},
		{
			name: "deployment added",
			revision: Revision{
				Prev: &Manifest{Content: []byte(`---
kind: Pod
metadata:
  name: foo
  namespace: bar
`)},
				Next: &Manifest{Content: []byte(`---
kind: Pod
metadata:
  name: foo
  namespace: bar
---
kind: Deployment
metadata:
  name: bar
`)},
			},
			expected: []ResourceSelector{},
		},
		{
			name: "multiple resources removed",
			revision: Revision{
				Prev: &Manifest{Content: []byte(`---
kind: Pod
metadata:
  name: foo
  namespace: bar
---
kind: ConfigMap
metadata:
  name: bar
  namespace: baz
`)},
				Next: &Manifest{Content: []byte(`---
---
kind: Deployment
metadata:
  name: bar
`)},
			},
			expected: []ResourceSelector{
				{
					Kind:      "Pod",
					Name:      "foo",
					Namespace: "bar",
				},
				{
					Kind:      "ConfigMap",
					Name:      "bar",
					Namespace: "baz",
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.revision.GetVanishedResources()

			assert.Equal(t, tc.expected, actual)
		})
	}
}
