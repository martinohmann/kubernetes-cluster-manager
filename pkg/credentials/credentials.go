package credentials

var (
	// Empty are empty credentials.
	Empty = Credentials{}
)

// Source provides credentials for a Kubernetes cluster.
type Source interface {
	// GetCredentials returns credentials for a Kubernetes cluster. Will return
	// an error if retrieving credentials fails.
	GetCredentials() (*Credentials, error)
}

// Credentials holds the credentials needed to communicate with a Kubernetes
// cluster.
type Credentials struct {
	Server     string `json:"server" yaml:"server"`
	Token      string `json:"token" yaml:"token"`
	Kubeconfig string `json:"kubeconfig" yaml:"kubeconfig"`
	Context    string `json:"context" yaml:"context"`
}

// Empty returns true if the credentials are empty.
func (c *Credentials) Empty() bool {
	return *c == Empty
}
