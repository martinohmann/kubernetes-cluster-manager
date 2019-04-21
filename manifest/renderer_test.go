package manifest

import (
	"testing"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/command"
	"github.com/stretchr/testify/assert"
)

func TestCreateRenderer(t *testing.T) {
	e := command.NewMockExecutor(nil)

	r, err := CreateRenderer("helm", &RendererOptions{}, e)
	assert.NoError(t, err)
	assert.IsType(t, &HelmRenderer{}, r)
}

func TestCreateRendererError(t *testing.T) {
	e := command.NewMockExecutor(nil)

	_, err := CreateRenderer("foo", &RendererOptions{}, e)
	assert.Error(t, err)
}
