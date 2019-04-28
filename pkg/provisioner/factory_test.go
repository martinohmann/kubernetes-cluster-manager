package provisioner

import (
	"testing"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kcm"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	r, err := Create("null", &kcm.ProvisionerOptions{}, nil)

	assert.NoError(t, err)
	assert.IsType(t, &Null{}, r)
}

func TestCreateError(t *testing.T) {
	_, err := Create("", &kcm.ProvisionerOptions{}, nil)

	assert.Error(t, err)
}
