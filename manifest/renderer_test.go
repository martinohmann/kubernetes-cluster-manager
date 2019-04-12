package manifest

import (
	"testing"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/command"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestCreateRenderer(t *testing.T) {
	e := command.NewMockExecutor(nil)
	cfg := &config.Config{ManifestRenderer: "helm"}

	r, err := CreateRenderer(cfg, e)
	assert.NoError(t, err)
	assert.IsType(t, &HelmRenderer{}, r)
}

func TestCreateRendererError(t *testing.T) {
	e := command.NewMockExecutor(nil)
	cfg := &config.Config{ManifestRenderer: "foo"}

	_, err := CreateRenderer(cfg, e)
	assert.Error(t, err)
}
