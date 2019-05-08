package kubernetes

import (
	"context"
	"errors"
	"testing"

	"github.com/martinohmann/kubernetes-cluster-manager/internal/commandtest"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/credentials"
	"github.com/stretchr/testify/assert"
)

func TestApplyManifest(t *testing.T) {
	commandtest.WithMockExecutor(func(executor commandtest.MockExecutor) {
		creds := &credentials.Credentials{
			Server: "https://localhost:6443",
			Token:  "sometoken",
		}

		kubectl := NewKubectl(creds)

		executor.ExpectCommand("kubectl apply -f - --server https://localhost:6443 --token sometoken")

		assert.NoError(t, kubectl.ApplyManifest(context.Background(), []byte{}))
		assert.NoError(t, executor.ExpectationsWereMet())
	})
}

func TestDeleteManifest(t *testing.T) {
	commandtest.WithMockExecutor(func(executor commandtest.MockExecutor) {
		creds := &credentials.Credentials{
			Kubeconfig: "/tmp/kubeconfig",
			Context:    "test",
		}

		kubectl := NewKubectl(creds)

		executor.ExpectCommand("kubectl delete -f - --ignore-not-found --context test --kubeconfig /tmp/kubeconfig")

		assert.NoError(t, kubectl.DeleteManifest(context.Background(), []byte{}))
		assert.NoError(t, executor.ExpectationsWereMet())
	})
}

func TestDeleteResource(t *testing.T) {
	commandtest.WithMockExecutor(func(executor commandtest.MockExecutor) {
		creds := &credentials.Credentials{
			Kubeconfig: "/tmp/kubeconfig",
		}

		selector := ResourceSelector{
			Name: "foo",
			Kind: "pod",
		}

		kubectl := NewKubectl(creds)

		executor.ExpectCommand("kubectl delete pod --ignore-not-found --namespace default --kubeconfig /tmp/kubeconfig foo")

		assert.NoError(t, kubectl.DeleteResource(context.Background(), selector))
		assert.NoError(t, executor.ExpectationsWereMet())
	})
}

func TestDeleteResources(t *testing.T) {
	commandtest.WithMockExecutor(func(executor commandtest.MockExecutor) {
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

		executor.ExpectCommand("kubectl delete pod .*")

		remaining, err := kubectl.DeleteResources(context.Background(), resources)

		assert.NoError(t, err)
		assert.Len(t, remaining, 0)

		assert.NoError(t, executor.ExpectationsWereMet())
	})
}

func TestDeleteResourcesError(t *testing.T) {
	commandtest.WithMockExecutor(func(executor commandtest.MockExecutor) {
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
				Kind: "deployment",
				Labels: map[string]string{
					"app.kubernetes.io/name": "bar",
				},
			},
		}

		kubectl := NewKubectl(creds)

		executor.ExpectCommand("kubectl delete pod .*")
		executor.ExpectCommand("kubectl delete deployment .*").WillReturnError(errors.New("error"))

		remaining, err := kubectl.DeleteResources(context.Background(), resources)

		assert.Error(t, err)
		assert.Len(t, remaining, 1)

		assert.NoError(t, executor.ExpectationsWereMet())
	})
}

func TestDeleteResourceLabels(t *testing.T) {
	commandtest.WithMockExecutor(func(executor commandtest.MockExecutor) {
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

		executor.ExpectCommand("kubectl delete pod --ignore-not-found --namespace default --kubeconfig /tmp/kubeconfig --selector app.kubernetes.io/name=foo,app.kubernetes.io/version=v0.0.1")

		assert.NoError(t, kubectl.DeleteResource(context.Background(), selector))
		assert.NoError(t, executor.ExpectationsWereMet())
	})
}

func TestDeleteResourceMissingSelector(t *testing.T) {
	commandtest.WithMockExecutor(func(executor commandtest.MockExecutor) {
		creds := &credentials.Credentials{
			Kubeconfig: "/tmp/kubeconfig",
		}

		selector := ResourceSelector{
			Kind: "pod",
		}

		kubectl := NewKubectl(creds)

		err := kubectl.DeleteResource(context.Background(), selector)

		assert.Error(t, err)
		assert.EqualError(t, err, "either a name or labels must be specified in the resource selector (kind=pod,namespace=default)")
		assert.NoError(t, executor.ExpectationsWereMet())
	})
}
