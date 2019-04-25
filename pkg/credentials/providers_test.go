package credentials

import (
	"testing"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kubernetes"
	"github.com/stretchr/testify/assert"
)

func TestStaticProvider(t *testing.T) {
	c := &kubernetes.Credentials{Kubeconfig: "~/.kube/config"}

	p := NewStaticProvider(c)

	credentials, err := p.GetCredentials()

	assert.NoError(t, err)
	assert.Exactly(t, c, credentials)
}
