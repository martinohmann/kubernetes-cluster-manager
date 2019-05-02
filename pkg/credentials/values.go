package credentials

import (
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kcm"
)

// ValueFetcherSource is a kcm.CredentialSource that retrieves Kubernetes
// credentials from an infrastructure provisioner.
type ValueFetcherSource struct {
	fetcher kcm.ValueFetcher
}

// NewValueFetcherSource creates a new ValueFetcherSource with given f.
func NewValueFetcherSource(f kcm.ValueFetcher) kcm.CredentialSource {
	return &ValueFetcherSource{f}
}

// GetCredentials implements kcm.CredentialSource.
func (s *ValueFetcherSource) GetCredentials() (*kcm.Credentials, error) {
	v, err := s.fetcher.Fetch()
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
