package manifest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestManifestFilename(t *testing.T) {
	m := &Manifest{Name: "manifest"}

	assert.Equal(t, "manifest.yaml", m.Filename())
}

func TestIsBlank(t *testing.T) {
	cases := []struct {
		name  string
		m     *Manifest
		blank bool
	}{
		{
			name:  "manifest nil",
			blank: true,
		},
		{
			name:  "manifest empty",
			m:     &Manifest{},
			blank: true,
		},
		{
			name:  "only whitespace",
			m:     &Manifest{Content: []byte("\n  \n ")},
			blank: true,
		},
		{
			name:  "whitespace and comments",
			m:     &Manifest{Content: []byte("# a comment \n  # another comment\n")},
			blank: true,
		},
		{
			name: "whitespace, comments and separators",
			m: &Manifest{Content: []byte(`---
# a comment

  # another comment
--- # inline comment

`)},
			blank: true,
		},
		{
			name: "whitespace, comments, separators and key-value pair",
			m: &Manifest{Content: []byte(`---
# a comment
somekey: somevalue
  # another comment
--- # inline comment

`)},
		},
		{
			name: "configmap",
			m: &Manifest{Content: []byte(`---
apiVersion: v1
kind: ConfigMap
metadata:
  name: kcm-chart
  labels:
    app.kubernetes.io/name: chart
    helm.sh/chart: cluster-0.1.0
    app.kubernetes.io/instance: kcm
data:
  SOMEVAR: someval

`)},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.blank, tc.m.IsBlank())
		})
	}
}

func TestManifestMatches(t *testing.T) {
	cases := []struct {
		name    string
		a, b    *Manifest
		matches bool
	}{
		{
			name:    "both nil",
			matches: true,
		},
		{
			name: "m nil",
			b:    &Manifest{},
		},
		{
			name: "other nil",
			a:    &Manifest{},
		},
		{
			name: "different name",
			a:    &Manifest{Name: "foo"},
			b:    &Manifest{Name: "bar"},
		},
		{
			name:    "same name",
			a:       &Manifest{Name: "foo"},
			b:       &Manifest{Name: "foo"},
			matches: true,
		},
		{
			name:    "same name, different content",
			a:       &Manifest{Name: "foo", Content: []byte(`a`)},
			b:       &Manifest{Name: "foo", Content: []byte(`b`)},
			matches: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.matches, tc.a.Matches(tc.b))
		})
	}
}

func TestReadDir(t *testing.T) {
	expected := []*Manifest{
		{
			Name: "foo",
			Content: []byte(`---
kind: Pod
metadata:
  name: foo
  namespace: bar
`),
		},
	}

	manifests, err := ReadDir("testdata/manifests")

	require.NoError(t, err)

	assert.Equal(t, expected, manifests)
}
