package kubernetes

import (
	"context"
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
