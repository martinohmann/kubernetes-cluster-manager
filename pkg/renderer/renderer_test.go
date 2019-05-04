package renderer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestManifestFilename(t *testing.T) {
	m := &Manifest{Name: "manifest"}

	assert.Equal(t, "manifest.yaml", m.Filename())
}
