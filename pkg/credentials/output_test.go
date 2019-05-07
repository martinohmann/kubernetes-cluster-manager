package credentials

import (
	"context"
	"testing"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kcm"
	"github.com/stretchr/testify/assert"
)

type testOutputter struct{}

func (testOutputter) Output(ctx context.Context) (kcm.Values, error) {
	return kcm.Values{
		"kubeconfig": "/tmp/kubeconfig",
	}, nil
}

func TestProvisionerOutputSource(t *testing.T) {
	p := NewProvisionerOutputSource(testOutputter{})

	credentials, err := p.GetCredentials(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, "/tmp/kubeconfig", credentials.Kubeconfig)
}
