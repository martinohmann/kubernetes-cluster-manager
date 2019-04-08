package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApplyDefaults(t *testing.T) {
	c := &Config{WorkingDir: "/tmp"}
	c.ApplyDefaults()

	assert.Equal(t, "/tmp/manifest.yaml", c.Manifest)
	assert.Equal(t, "/tmp/deletions.yaml", c.Deletions)
	assert.Equal(t, "/tmp/values.yaml", c.Values)
	assert.Equal(t, "/tmp/cluster", c.Helm.Chart)
}
