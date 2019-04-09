package provisioner

import (
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/api"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/config"
)

// updateCredentialsFromValues tries to update the cluster credentials in the
// config from values retrieved from the infrastructure manager. It will not
// overwrite config values that are already set.
func updateCredentialsFromValues(cfg *config.Config, values api.Values) {
	if s, ok := values["server"].(string); ok && cfg.Server == "" {
		cfg.Server = s
	}

	if t, ok := values["token"].(string); ok && cfg.Token == "" {
		cfg.Token = t
	}

	if k, ok := values["kubeconfig"].(string); ok && cfg.Kubeconfig == "" {
		cfg.Kubeconfig = k
	}
}
