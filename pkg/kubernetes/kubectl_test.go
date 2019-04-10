package kubernetes

import (
	"testing"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/api"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/command"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestApplyManifest(t *testing.T) {
	executor := command.NewMockExecutor()
	cfg := &config.ClusterConfig{
		Server: "https://localhost:6443",
		Token:  "sometoken",
	}

	kubectl := NewKubectl(cfg, executor)

	err := kubectl.ApplyManifest(api.Manifest{})

	assert.NoError(t, err)
	if assert.Len(t, executor.ExecutedCommands, 1) {
		assert.Equal(
			t,
			executor.ExecutedCommands[0],
			"kubectl apply -f - --server https://localhost:6443 --token sometoken",
		)
	}
}

func TestDeleteManifest(t *testing.T) {
	executor := command.NewMockExecutor()
	cfg := &config.ClusterConfig{
		Kubeconfig: "/tmp/kubeconfig",
	}

	kubectl := NewKubectl(cfg, executor)

	err := kubectl.DeleteManifest(api.Manifest{})

	assert.NoError(t, err)
	if assert.Len(t, executor.ExecutedCommands, 1) {
		assert.Equal(
			t,
			executor.ExecutedCommands[0],
			"kubectl delete -f - --ignore-not-found --kubeconfig /tmp/kubeconfig",
		)
	}
}

func TestDeleteResource(t *testing.T) {
	executor := command.NewMockExecutor()
	cfg := &config.ClusterConfig{
		Kubeconfig: "/tmp/kubeconfig",
	}

	resource := &api.Deletion{
		Name: "foo",
		Kind: "pod",
	}

	kubectl := NewKubectl(cfg, executor)

	err := kubectl.DeleteResource(resource)

	assert.NoError(t, err)
	if assert.Len(t, executor.ExecutedCommands, 1) {
		assert.Equal(
			t,
			executor.ExecutedCommands[0],
			"kubectl delete pod --ignore-not-found --namespace default --kubeconfig /tmp/kubeconfig foo",
		)
	}
}

func TestDeleteResourceLabels(t *testing.T) {
	executor := command.NewMockExecutor()
	cfg := &config.ClusterConfig{
		Kubeconfig: "/tmp/kubeconfig",
	}

	resource := &api.Deletion{
		Kind: "pod",
		Labels: map[string]string{
			"app.kubernetes.io/name":    "foo",
			"app.kubernetes.io/version": "v0.0.1",
		},
	}

	kubectl := NewKubectl(cfg, executor)

	err := kubectl.DeleteResource(resource)

	assert.NoError(t, err)
	if assert.Len(t, executor.ExecutedCommands, 1) {
		assert.Equal(
			t,
			executor.ExecutedCommands[0],
			"kubectl delete pod --ignore-not-found --namespace default --kubeconfig /tmp/kubeconfig --selector app.kubernetes.io/name=foo,app.kubernetes.io/version=v0.0.1",
		)
	}
}

func TestDeleteResourceMissingSelector(t *testing.T) {
	executor := command.NewMockExecutor()
	cfg := &config.ClusterConfig{
		Kubeconfig: "/tmp/kubeconfig",
	}

	resource := &api.Deletion{
		Kind: "pod",
	}

	kubectl := NewKubectl(cfg, executor)

	err := kubectl.DeleteResource(resource)

	assert.Error(t, err)
	assert.EqualError(t, err, "either a name or labels must be specified for a deletion (kind=pod,namespace=default)")
}
