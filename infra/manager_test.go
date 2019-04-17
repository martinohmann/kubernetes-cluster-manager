package infra

import (
	"testing"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/command"
	"github.com/stretchr/testify/assert"
)

func TestCreateManager(t *testing.T) {
	e := command.NewMockExecutor(nil)

	r, err := CreateManager("terraform", &ManagerOptions{}, e)
	assert.NoError(t, err)
	assert.IsType(t, &TerraformManager{}, r)
}

func TestCreateManagerError(t *testing.T) {
	e := command.NewMockExecutor(nil)

	_, err := CreateManager("foo", &ManagerOptions{}, e)
	assert.Error(t, err)
}
