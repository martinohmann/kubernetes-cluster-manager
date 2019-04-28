package provisioner

import (
	"testing"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/command"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kcm"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/assert"
)

func TestMinikubeProvision(t *testing.T) {
	executor := command.NewMockExecutor(nil)

	m := NewMinikube(executor)

	executor.Command("minikube status").WillError()
	executor.Command("minikube start").WillSucceed()

	err := m.Provision()

	assert.NoError(t, err)
}

func TestMinikubeFetch(t *testing.T) {
	executor := command.NewMockExecutor(nil)

	m := NewMinikube(executor)

	output := `127.0.0.1`

	executor.Command("minikube ip").WillReturn(output)

	home, _ := homedir.Dir()

	expectedValues := kcm.Values{
		"context":    "minikube",
		"kubeconfig": home + "/.kube/config",
		"server":     "https://127.0.0.1:8443",
	}

	values, err := m.Fetch()

	assert.NoError(t, err)
	assert.Equal(t, expectedValues, values)
}

func TestMinikubeDestroy(t *testing.T) {
	executor := command.NewMockExecutor(nil)

	m := NewMinikube(executor)

	executor.Command("minikube status").WillSucceed()
	executor.Command("minikube delete").WillSucceed()

	err := m.Destroy()

	assert.NoError(t, err)
}
