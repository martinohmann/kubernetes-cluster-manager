package infra

import (
	"testing"

	"github.com/martinohmann/kubernetes-cluster-manager/infra/terraform"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/command"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestCreateManager(t *testing.T) {
	e := command.NewMockExecutor(nil)
	cfg := &config.Config{InfraManager: "terraform"}

	r, err := CreateManager(cfg, e)
	assert.NoError(t, err)
	assert.IsType(t, &terraform.InfraManager{}, r)
}

func TestCreateManagerError(t *testing.T) {
	e := command.NewMockExecutor(nil)
	cfg := &config.Config{InfraManager: "foo"}

	_, err := CreateManager(cfg, e)
	assert.Error(t, err)
}
