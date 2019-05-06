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
