package renderer

import (
	"testing"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kcm"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	r, err := Create("helm", &kcm.RendererOptions{})

	assert.NoError(t, err)
	assert.IsType(t, &Helm{}, r)
}

func TestCreateError(t *testing.T) {
	_, err := Create("", &kcm.RendererOptions{})

	assert.Error(t, err)
}