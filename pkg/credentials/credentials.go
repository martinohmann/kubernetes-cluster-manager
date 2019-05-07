package credentials

import "context"

var (
	// Empty are empty credentials.
	Empty = Credentials{}
)

// Source provides credentials for a Kubernetes cluster.
type Source interface {
	// GetCredentials returns credentials for a Kubernetes cluster. Will return
	// an error if retrieving credentials fails.
	GetCredentials(context.Context) (*Credentials, error)
}

// Credentials holds the credentials needed to communicate with a Kubernetes
// cluster.
type Credentials struct {
	Server     string `json:"server,omitempty" yaml:"server,omitempty"`
	Token      string `json:"token,omitempty" yaml:"token,omitempty"`
	Kubeconfig string `json:"kubeconfig,omitempty" yaml:"kubeconfig,omitempty"`
	Context    string `json:"context,omitempty" yaml:"context,omitempty"`
}

// Empty returns true if the credentials are empty.
func (c *Credentials) Empty() bool {
	return *c == Empty
}
