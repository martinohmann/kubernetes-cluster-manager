package resource

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSlice_Bytes(t *testing.T) {
	s := Slice{
		{Content: []byte(`apiVersion: v1
kind: ConfigMap
metadata:
  name: bar
  namespace: baz
`)},
		{Content: []byte(`apiVersion: v1
kind: Pod
metadata:
  name: foo
  namespace: bar
`)},
	}

	expected := []byte(`---
apiVersion: v1
kind: ConfigMap
metadata:
  name: bar
  namespace: baz

---
apiVersion: v1
kind: Pod
metadata:
  name: foo
  namespace: bar

`)

	assert.Equal(t, string(expected), string(s.Bytes()))
}

func TestSlice_String(t *testing.T) {
	s := Slice{
		{Name: "foo", Kind: "Pod"},
		{Name: "bar", Kind: "Deployment"},
		{Name: "baz", Kind: "StatefulSet"},
		{Name: "prometheus", Kind: "CustomResourceDefinition"},
	}

	expected := `pod/foo
deployment/bar
statefulset/baz
customresourcedefinition/prometheus`

	assert.Equal(t, expected, s.String())
}

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

	r, _ := New(nil, Head{Kind: "Pod", Metadata: Metadata{Name: "foo", Namespace: "kube-system"}})

	resources := []*Resource{
		nil,
		{Name: "foo", Kind: "Pod", Namespace: "kube-system"},
		r,
	}

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
