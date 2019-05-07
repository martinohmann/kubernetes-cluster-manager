package provisioner

import (
	"context"
	"testing"

	"github.com/martinohmann/kubernetes-cluster-manager/internal/commandtest"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kcm"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/assert"
)

func TestMinikubeProvision(t *testing.T) {
	commandtest.WithMockExecutor(func(executor *commandtest.MockExecutor) {
		m := NewMinikube(&Options{})

		executor.Command("minikube status").WillError()
		executor.Command("minikube start").WillSucceed()

		err := m.Provision(context.Background())

		assert.NoError(t, err)
	})
}

func TestMinikubeOutput(t *testing.T) {
	m := &Minikube{}

	home, _ := homedir.Dir()

	expectedValues := kcm.Values{
		"context":    "minikube",
		"kubeconfig": home + "/.kube/config",
	}

	values, err := m.Output(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, expectedValues, values)
}

func TestMinikubeDestroy(t *testing.T) {
	commandtest.WithMockExecutor(func(executor *commandtest.MockExecutor) {
		m := &Minikube{}

		executor.Command("minikube status").WillSucceed()
		executor.Command("minikube delete").WillSucceed()

		err := m.Destroy(context.Background())

		assert.NoError(t, err)
	})
}
