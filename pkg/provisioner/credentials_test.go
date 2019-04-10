package provisioner

import (
	"testing"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/api"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestUpdateCredentialsFromValues(t *testing.T) {
	cfg := config.ClusterConfig{
		Kubeconfig: "~/.kube/config",
		Token:      "supersecret",
	}

	values := api.Values{
		"server":     "https://localhost:6443",
		"kubeconfig": "/tmp/kubeconfig",
	}

	updateCredentialsFromValues(&cfg, values)

	assert.Equal(t, "https://localhost:6443", cfg.Server)
	assert.Equal(t, "~/.kube/config", cfg.Kubeconfig)
	assert.Equal(t, "supersecret", cfg.Token)
}
