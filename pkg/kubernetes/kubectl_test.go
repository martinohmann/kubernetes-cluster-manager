package kubernetes

import (
	"testing"

	"github.com/martinohmann/kubernetes-cluster-manager/internal/commandtest"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/credentials"
	"github.com/stretchr/testify/assert"
)

func TestApplyManifest(t *testing.T) {
	commandtest.WithMockExecutor(func(executor *commandtest.MockExecutor) {
		creds := &credentials.Credentials{
			Server: "https://localhost:6443",
			Token:  "sometoken",
		}

		kubectl := NewKubectl(creds)

		err := kubectl.ApplyManifest([]byte{})

		assert.NoError(t, err)
		if assert.Len(t, executor.ExecutedCommands, 1) {
			assert.Equal(
				t,
				executor.ExecutedCommands[0],
				"kubectl apply -f - --server https://localhost:6443 --token sometoken",
			)
		}
	})
}

func TestDeleteManifest(t *testing.T) {
	commandtest.WithMockExecutor(func(executor *commandtest.MockExecutor) {
		creds := &credentials.Credentials{
			Kubeconfig: "/tmp/kubeconfig",
			Context:    "test",
		}

		kubectl := NewKubectl(creds)

		err := kubectl.DeleteManifest([]byte{})

		assert.NoError(t, err)
		if assert.Len(t, executor.ExecutedCommands, 1) {
			assert.Equal(
				t,
				executor.ExecutedCommands[0],
				"kubectl delete -f - --ignore-not-found --context test --kubeconfig /tmp/kubeconfig",
			)
		}
	})
}

func TestDeleteResource(t *testing.T) {
	commandtest.WithMockExecutor(func(executor *commandtest.MockExecutor) {
		creds := &credentials.Credentials{
			Kubeconfig: "/tmp/kubeconfig",
		}

		selector := ResourceSelector{
			Name: "foo",
			Kind: "pod",
		}

		kubectl := NewKubectl(creds)

		err := kubectl.DeleteResource(selector)

		assert.NoError(t, err)
		if assert.Len(t, executor.ExecutedCommands, 1) {
			assert.Equal(
				t,
				executor.ExecutedCommands[0],
				"kubectl delete pod --ignore-not-found --namespace default --kubeconfig /tmp/kubeconfig foo",
			)
		}
	})
}

func TestDeleteResources(t *testing.T) {
	commandtest.WithMockExecutor(func(executor *commandtest.MockExecutor) {
		creds := &credentials.Credentials{
			Kubeconfig: "/tmp/kubeconfig",
		}

		resources := []ResourceSelector{
			{
				Kind: "pod",
				Labels: map[string]string{
					"app.kubernetes.io/name": "foo",
				},
			},
		}

		kubectl := NewKubectl(creds)

		executor.NextCommand().WillSucceed()

		remaining, err := kubectl.DeleteResources(resources)

		assert.NoError(t, err)
		assert.Len(t, remaining, 0)
	})
}

func TestDeleteResourcesError(t *testing.T) {
	commandtest.WithMockExecutor(func(executor *commandtest.MockExecutor) {
		creds := &credentials.Credentials{
			Kubeconfig: "/tmp/kubeconfig",
		}

		resources := []ResourceSelector{
			{
				Kind: "pod",
				Labels: map[string]string{
					"app.kubernetes.io/name": "foo",
				},
			},
			{
				Kind: "bar",
				Labels: map[string]string{
					"app.kubernetes.io/name": "bar",
				},
			},
		}

		kubectl := NewKubectl(creds)

		executor.NextCommand().WillSucceed()
		executor.NextCommand().WillError()

		remaining, err := kubectl.DeleteResources(resources)

		assert.Error(t, err)
		assert.Len(t, remaining, 1)
	})
}

func TestDeleteResourceLabels(t *testing.T) {
	commandtest.WithMockExecutor(func(executor *commandtest.MockExecutor) {
		creds := &credentials.Credentials{
			Kubeconfig: "/tmp/kubeconfig",
		}

		selector := ResourceSelector{
			Kind: "pod",
			Labels: map[string]string{
				"app.kubernetes.io/name":    "foo",
				"app.kubernetes.io/version": "v0.0.1",
			},
		}

		kubectl := NewKubectl(creds)

		err := kubectl.DeleteResource(selector)

		assert.NoError(t, err)
		if assert.Len(t, executor.ExecutedCommands, 1) {
			assert.Equal(
				t,
				executor.ExecutedCommands[0],
				"kubectl delete pod --ignore-not-found --namespace default --kubeconfig /tmp/kubeconfig --selector app.kubernetes.io/name=foo,app.kubernetes.io/version=v0.0.1",
			)
		}
	})
}

func TestDeleteResourceMissingSelector(t *testing.T) {
	commandtest.WithMockExecutor(func(executor *commandtest.MockExecutor) {
		creds := &credentials.Credentials{
			Kubeconfig: "/tmp/kubeconfig",
		}

		selector := ResourceSelector{
			Kind: "pod",
		}

		kubectl := NewKubectl(creds)

		err := kubectl.DeleteResource(selector)

		assert.Error(t, err)
		assert.EqualError(t, err, "either a name or labels must be specified in the resource selector (kind=pod,namespace=default)")
	})
}
