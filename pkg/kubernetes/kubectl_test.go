package kubernetes

import (
	"context"
	"errors"
	"testing"

	"github.com/martinohmann/kubernetes-cluster-manager/internal/commandtest"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/credentials"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/resource"
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
			Context:    "test",
		}

		kubectl := NewKubectl(creds)

		executor.ExpectCommand("kubectl delete statefulset foo --namespace default --ignore-not-found --context test --kubeconfig /tmp/kubeconfig")

		res := resource.Head{
			Kind: resource.StatefulSet,
			Metadata: resource.Metadata{
				Name: "foo",
			},
		}

		assert.NoError(t, kubectl.DeleteResource(context.Background(), res))
		assert.NoError(t, executor.ExpectationsWereMet())
	})
}

func TestValidationErrors(t *testing.T) {
	commandtest.WithMockExecutor(func(executor commandtest.MockExecutor) {
		kubectl := NewKubectl(&credentials.Credentials{})

		executor.ExpectCommand("kubectl delete someunknownkind foo --namespace default --ignore-not-found").
			WillReturnError(errors.New(`error: the server doesn't have a resource type "someunknownkind"`))

		res := resource.Head{
			Kind: "SomeUnknownKind",
			Metadata: resource.Metadata{
				Name: "foo",
			},
		}

		assert.Error(t, kubectl.DeleteResource(context.Background(), res))
		assert.NoError(t, executor.ExpectationsWereMet())
	})
}

func TestClusterInfo(t *testing.T) {
	commandtest.WithMockExecutor(func(executor commandtest.MockExecutor) {
		creds := &credentials.Credentials{
			Kubeconfig: "/tmp/kubeconfig",
			Context:    "test",
		}

		kubectl := NewKubectl(creds)

		executor.ExpectCommand("kubectl cluster-info --context test --kubeconfig /tmp/kubeconfig")

		_, err := kubectl.ClusterInfo(context.Background())

		assert.NoError(t, err)
		assert.NoError(t, executor.ExpectationsWereMet())
	})
}
