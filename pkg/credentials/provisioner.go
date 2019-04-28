package credentials

import (
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kcm"
)

// ProvisionerSource is a Source that retrieves Kubernetes credentials from an
// infrastructure provisioner.
type ProvisionerSource struct {
	provisioner kcm.Provisioner
}

// NewInfraSource creates a new InfraSource with given manager.
func NewProvisionerSource(p kcm.Provisioner) kcm.CredentialSource {
	return &ProvisionerSource{p}
}

// GetCredentials implements Source.
func (p *ProvisionerSource) GetCredentials() (*kcm.Credentials, error) {
	v, err := p.provisioner.Fetch()
	if err != nil {
		return nil, err
	}

	c := &kcm.Credentials{}

	if server, ok := v["server"].(string); ok {
		c.Server = server
	}

	if token, ok := v["token"].(string); ok {
		c.Token = token
	}

	if kubeconfig, ok := v["kubeconfig"].(string); ok {
		c.Kubeconfig = kubeconfig
	}

	if context, ok := v["context"].(string); ok {
		c.Context = context
	}

	return c, nil
}
