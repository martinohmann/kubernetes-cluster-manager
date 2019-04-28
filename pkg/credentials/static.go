package credentials

import (
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kcm"
)

// StaticCredentials is a Source that holds static Kubernetes credentials.
type StaticCredentials struct {
	c *kcm.Credentials
}

// NewStaticCredentials creates a new StaticCredentials source with given credentials.
func NewStaticCredentials(c *kcm.Credentials) kcm.CredentialSource {
	return &StaticCredentials{c}
}

// GetCredentials implements Source.
func (p *StaticCredentials) GetCredentials() (*kcm.Credentials, error) {
	return p.c, nil
}
