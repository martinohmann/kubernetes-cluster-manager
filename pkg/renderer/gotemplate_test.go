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

	expectedOneContent := []byte(`---
# Source: one/test.yaml
    BAZ

`)

	expectedTwoContent := []byte(`---
# Source: two/00-test.yaml
---
blah

---
# Source: two/test.yaml
qux

`)

	assert.Equal(t, "one", manifests[0].Name)
	assert.Equal(t, expectedOneContent, manifests[0].Content)
	assert.Equal(t, "two", manifests[1].Name)
	assert.Equal(t, expectedTwoContent, manifests[1].Content)
}
