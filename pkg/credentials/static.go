package credentials

// StaticSource is a Source that holds static Kubernetes credentials.
type StaticSource struct {
	c *Credentials
}

// NewStaticSource creates a new StaticSource source with given credentials.
func NewStaticSource(c *Credentials) Source {
	return &StaticSource{c}
}

// GetCredentials implements Source.
func (p *StaticSource) GetCredentials() (*Credentials, error) {
	return p.c, nil
}
