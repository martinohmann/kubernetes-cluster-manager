package credentials

import (
	"testing"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kcm"
	"github.com/stretchr/testify/assert"
)

func TestStaticCredentials(t *testing.T) {
	c := &kcm.Credentials{Kubeconfig: "~/.kube/config"}

	p := NewStaticCredentials(c)

	credentials, err := p.GetCredentials()

	assert.NoError(t, err)
	assert.Exactly(t, c, credentials)
}
