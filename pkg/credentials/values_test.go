package credentials

import (
	"testing"

	"github.com/martinohmann/kubernetes-cluster-manager/internal/commandtest"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/command"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kcm"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/provisioner"
	"github.com/stretchr/testify/assert"
)

func TestValueFetcherSource(t *testing.T) {
	commandtest.WithMockExecutor(func(executor *commandtest.MockExecutor) {
		executor.Command("terraform output --json").WillReturn(`{"kubeconfig":{"value": "/tmp/kubeconfig"}}`)

		p := NewValueFetcherSource(provisioner.NewTerraform(&kcm.TerraformOptions{}))

		credentials, err := p.GetCredentials()

		assert.NoError(t, err)
		assert.Equal(t, "/tmp/kubeconfig", credentials.Kubeconfig)
	}, command.NewExecutor(nil))
}
