package renderer

import (
	"testing"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kcm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGoTemplateRenderManifests(t *testing.T) {
	o := &Options{
		TemplatesDir: "testdata/gotemplate",
	}

	v := kcm.Values{
		"foo": "bar",
		"bar": "baz",
		"baz": "qux",
	}

	r := NewGoTemplate(o)

	manifests, err := r.RenderManifests(v)

	require.NoError(t, err)
	require.Len(t, manifests, 2)

	assert.Equal(t, "one.yaml", manifests[0].Filename)
	assert.Equal(t, "---\n    BAZ\n", string(manifests[0].Content))
	assert.Equal(t, "two.yaml", manifests[1].Filename)
	assert.Equal(t, "---\n---\nblah\n---\nqux\n", string(manifests[1].Content))
}
