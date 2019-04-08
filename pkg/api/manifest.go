package api

// Manifest defines a Kubernetes manifest.
type Manifest struct {
	Content []byte `json:"content" yaml:"content"`
}

func NewManifestFromString(manifest string) *Manifest {
	return &Manifest{Content: []byte(manifest)}
}

// String implements fmt.Stringer.
func (m *Manifest) String() string {
	return string(m.Content)
}
