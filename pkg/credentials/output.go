package credentials

import (
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/provisioner"
)

// ProvisionerOutputSource is a Source that retrieves Kubernetes
// credentials from an infrastructure provisioner.
type ProvisionerOutputSource struct {
	o provisioner.Outputter
}

// NewProvisionerOutputSource creates a new ProvisionerOutputSource with given o.
func NewProvisionerOutputSource(o provisioner.Outputter) Source {
	return &ProvisionerOutputSource{o}
}

// GetCredentials implements Source.
func (s *ProvisionerOutputSource) GetCredentials() (*Credentials, error) {
	v, err := s.o.Output()
	if err != nil {
		return nil, err
	}

	c := &Credentials{}

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
