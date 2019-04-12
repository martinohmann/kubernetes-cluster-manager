package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApplyDefaults(t *testing.T) {
	c := &Config{WorkingDir: "/tmp"}
	c.ApplyDefaults()

	assert.Equal(t, "/tmp/manifest.yaml", c.Manifest)
	assert.Equal(t, "/tmp/deletions.yaml", c.Deletions)
	assert.Equal(t, "/tmp/values.yaml", c.Values)
	assert.Equal(t, "/tmp/cluster", c.Helm.Chart)
}

func TestUpdateClusterConfig(t *testing.T) {
	cfg := ClusterConfig{
		Kubeconfig: "~/.kube/config",
		Token:      "supersecret",
	}

	values := map[string]interface{}{
		"server":     "https://localhost:6443",
		"kubeconfig": "/tmp/kubeconfig",
	}

	cfg.Update(values)

	assert.Equal(t, "https://localhost:6443", cfg.Server)
	assert.Equal(t, "~/.kube/config", cfg.Kubeconfig)
	assert.Equal(t, "supersecret", cfg.Token)
}
