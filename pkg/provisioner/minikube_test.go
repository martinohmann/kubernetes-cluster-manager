package provisioner

import (
	"context"
	"testing"

	"github.com/martinohmann/kubernetes-cluster-manager/internal/commandtest"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestMinikubeProvision(t *testing.T) {
	commandtest.WithMockExecutor(func(executor commandtest.MockExecutor) {
		m := NewMinikube(&Options{})

		executor.ExpectCommand("minikube status").WillReturnError(errors.New("not running"))
		executor.ExpectCommand("minikube start")

		assert.NoError(t, m.Provision(context.Background()))
		assert.NoError(t, executor.ExpectationsWereMet())
	})
}

func TestMinikubeOutput(t *testing.T) {
	m := &Minikube{}

	home, _ := homedir.Dir()

	expectedValues := map[string]interface{}{
		"context":    "minikube",
		"kubeconfig": home + "/.kube/config",
	}

	values, err := m.Output(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, expectedValues, values)
}

func TestMinikubeDestroy(t *testing.T) {
	commandtest.WithMockExecutor(func(executor commandtest.MockExecutor) {
		m := &Minikube{}

		executor.ExpectCommand("minikube status")
		executor.ExpectCommand("minikube delete")

		assert.NoError(t, m.Destroy(context.Background()))
		assert.NoError(t, executor.ExpectationsWereMet())
	})
}
