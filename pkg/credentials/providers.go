package credentials

import (
	"github.com/martinohmann/kubernetes-cluster-manager/infra"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kubernetes"
)

// Provider provides credentials for a Kubernetes cluster.
type Provider interface {
	// GetCredentials returns credentials for a Kubernetes cluster. Will return
	// an error if retrieving credentials fails.
	GetCredentials() (*kubernetes.Credentials, error)
}

// StaticProvider is a Provider that holds static Kubernetes credentials.
type StaticProvider struct {
	c *kubernetes.Credentials
}

// NewStaticProvider creates a new StaticProvider with given credentials.
func NewStaticProvider(c *kubernetes.Credentials) Provider {
	return &StaticProvider{c}
}

// GetCredentials implements Provider.
func (p *StaticProvider) GetCredentials() (*kubernetes.Credentials, error) {
	return p.c, nil
}

// InfraProvider is a Provider that retrieves Kubernetes credentials from an
// infrastructure manager.
type InfraProvider struct {
	m infra.Manager
}

// NewInfraProvider creates a new InfraProvider with given manager.
func NewInfraProvider(m infra.Manager) Provider {
	return &InfraProvider{m}
}

// GetCredentials implements Provider.
func (p *InfraProvider) GetCredentials() (*kubernetes.Credentials, error) {
	v, err := p.m.GetValues()
	if err != nil {
		return nil, err
	}

	c := &kubernetes.Credentials{}

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
