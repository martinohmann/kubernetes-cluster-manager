package credentials

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStaticSource(t *testing.T) {
	c := &Credentials{Kubeconfig: "~/.kube/config"}

	p := NewStaticSource(c)

	credentials, err := p.GetCredentials(context.Background())

	assert.NoError(t, err)
	assert.Exactly(t, c, credentials)
}
