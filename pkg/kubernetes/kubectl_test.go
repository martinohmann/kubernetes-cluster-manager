package kubernetes

import (
	"testing"

	"github.com/martinohmann/kubernetes-cluster-manager/internal/commandtest"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kcm"
	"github.com/stretchr/testify/assert"
)

func TestApplyManifest(t *testing.T) {
	commandtest.WithMockExecutor(func(executor *commandtest.MockExecutor) {
		creds := &kcm.Credentials{
			Server: "https://localhost:6443",
			Token:  "sometoken",
		}

		kubectl := NewKubectl(creds)

		err := kubectl.ApplyManifest(&kcm.Manifest{})

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
		creds := &kcm.Credentials{
			Kubeconfig: "/tmp/kubeconfig",
			Context:    "test",
		}

		kubectl := NewKubectl(creds)

		err := kubectl.DeleteManifest(&kcm.Manifest{})

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
		creds := &kcm.Credentials{
			Kubeconfig: "/tmp/kubeconfig",
		}

		resource := &kcm.Deletion{
			Name: "foo",
			Kind: "pod",
		}

		kubectl := NewKubectl(creds)

		err := kubectl.DeleteResource(resource)

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

func TestDeleteResourceLabels(t *testing.T) {
	commandtest.WithMockExecutor(func(executor *commandtest.MockExecutor) {
		creds := &kcm.Credentials{
			Kubeconfig: "/tmp/kubeconfig",
		}

		resource := &kcm.Deletion{
			Kind: "pod",
			Labels: map[string]string{
				"app.kubernetes.io/name":    "foo",
				"app.kubernetes.io/version": "v0.0.1",
			},
		}

		kubectl := NewKubectl(creds)

		err := kubectl.DeleteResource(resource)

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
		creds := &kcm.Credentials{
			Kubeconfig: "/tmp/kubeconfig",
		}

		resource := &kcm.Deletion{
			Kind: "pod",
		}

		kubectl := NewKubectl(creds)

		err := kubectl.DeleteResource(resource)

		assert.Error(t, err)
		assert.EqualError(t, err, "either a name or labels must be specified for a deletion (kind=pod,namespace=default)")
	})
}
