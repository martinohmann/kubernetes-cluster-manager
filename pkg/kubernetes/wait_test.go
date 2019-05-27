package kubernetes

import (
	"context"
	"testing"
	"time"

	"github.com/martinohmann/kubernetes-cluster-manager/internal/commandtest"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/credentials"
	"github.com/stretchr/testify/assert"
)

func TestWait(t *testing.T) {
	commandtest.WithMockExecutor(func(executor commandtest.MockExecutor) {
		creds := &credentials.Credentials{
			Kubeconfig: "/tmp/kubeconfig",
			Context:    "test",
		}

		kubectl := NewKubectl(creds)

		executor.ExpectCommand("kubectl wait --for condition=complete --namespace bar job/foo --timeout 600s --context test --kubeconfig /tmp/kubeconfig")

		opts := WaitOptions{
			Kind:      "job",
			Name:      "foo",
			Namespace: "bar",
			For:       "condition=complete",
			Timeout:   10 * time.Minute,
		}

		assert.NoError(t, kubectl.Wait(context.Background(), opts))
		assert.NoError(t, executor.ExpectationsWereMet())
	})
}
