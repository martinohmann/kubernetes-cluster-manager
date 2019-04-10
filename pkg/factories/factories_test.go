package factories

import (
	"testing"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/command"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/config"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kubernetes/helm"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/terraform"
	"github.com/stretchr/testify/assert"
)

func TestCreateManifestRenderer(t *testing.T) {
	e := command.NewMockExecutor()
	cfg := &config.Config{ManifestRenderer: "helm"}

	r, err := CreateManifestRenderer(cfg, e)
	assert.NoError(t, err)
	assert.IsType(t, &helm.ManifestRenderer{}, r)
}

func TestCreateManifestRendererError(t *testing.T) {
	e := command.NewMockExecutor()
	cfg := &config.Config{ManifestRenderer: "foo"}

	_, err := CreateManifestRenderer(cfg, e)
	assert.Error(t, err)
}

func TestCreateInfraManager(t *testing.T) {
	e := command.NewMockExecutor()
	cfg := &config.Config{InfraManager: "terraform"}

	r, err := CreateInfraManager(cfg, e)
	assert.NoError(t, err)
	assert.IsType(t, &terraform.InfraManager{}, r)
}

func TestCreateInfraManagerError(t *testing.T) {
	e := command.NewMockExecutor()
	cfg := &config.Config{InfraManager: "foo"}

	_, err := CreateInfraManager(cfg, e)
	assert.Error(t, err)
}
