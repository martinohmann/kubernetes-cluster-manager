package kubernetes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClusterOptions(t *testing.T) {
	o := &ClusterOptions{
		Server: "https://localhost:6443",
		Token:  "secret",
	}

	v := map[string]interface{}{
		"server":     "https://localhost:8443",
		"kubeconfig": "~/.kube/config",
		"somekey":    "somevalues",
	}

	o.Update(v)

	assert.Equal(t, "https://localhost:6443", o.Server)
	assert.Equal(t, "secret", o.Token)
	assert.Equal(t, "~/.kube/config", o.Kubeconfig)
	assert.Equal(t, "", o.Context)
}
