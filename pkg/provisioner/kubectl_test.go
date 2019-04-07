package provisioner

import (
	"testing"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/api"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/command"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestApplyManifest(t *testing.T) {
	executor := command.NewMockExecutor()
	cfg := &config.Config{
		DryRun: true,
		Server: "https://localhost:6443",
		Token:  "sometoken",
	}

	kubectl := NewKubectl(cfg, executor)

	err := kubectl.ApplyManifest(&api.Manifest{})

	assert.NoError(t, err)
	if assert.Len(t, executor.ExecutedCommands, 1) {
		assert.Equal(
			t,
			executor.ExecutedCommands[0],
			"kubectl apply -f - --server https://localhost:6443 --token sometoken --dry-run",
		)
	}
}

func TestDeleteManifest(t *testing.T) {
	executor := command.NewMockExecutor()
	cfg := &config.Config{
		Kubeconfig: "/tmp/kubeconfig",
	}

	kubectl := NewKubectl(cfg, executor)

	err := kubectl.DeleteManifest(&api.Manifest{})

	assert.NoError(t, err)
	if assert.Len(t, executor.ExecutedCommands, 1) {
		assert.Equal(
			t,
			executor.ExecutedCommands[0],
			"kubectl delete -f - --ignore-not-found --kubeconfig /tmp/kubeconfig",
		)
	}
}

func TestDeleteManifestDryRun(t *testing.T) {
	executor := command.NewMockExecutor()
	cfg := &config.Config{DryRun: true}
	manifest := &api.Manifest{Content: []byte("---")}

	kubectl := NewKubectl(cfg, executor)

	err := kubectl.DeleteManifest(manifest)

	assert.NoError(t, err)
	assert.Len(t, executor.ExecutedCommands, 0)
}

func TestDeleteResource(t *testing.T) {
	executor := command.NewMockExecutor()
	cfg := &config.Config{
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

func TestDeleteResourceDryRun(t *testing.T) {
	executor := command.NewMockExecutor()
	cfg := &config.Config{
		DryRun: true,
	}

	kubectl := NewKubectl(cfg, executor)

	err := kubectl.DeleteResource(&api.Deletion{})

	assert.NoError(t, err)
	assert.Len(t, executor.ExecutedCommands, 0)
}
