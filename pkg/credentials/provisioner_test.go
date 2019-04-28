package credentials

import (
	"testing"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/provisioner"
	"github.com/stretchr/testify/assert"
)

func TestProvisionerSource(t *testing.T) {
	p := NewProvisionerSource(&provisioner.Null{})

	credentials, err := p.GetCredentials()

	assert.NoError(t, err)
	assert.NotNil(t, credentials)
}
