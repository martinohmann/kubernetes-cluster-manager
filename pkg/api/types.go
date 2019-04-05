package api

type Deletions struct {
	PreApply   []Deletion `json:"preApply" yaml:"preApply"`
	PostApply  []Deletion `json:"postApply" yaml:"postApply"`
	PreDestroy []Deletion `json:"preDestroy" yaml:"preDestroy"`
}

type Deletion struct {
	Kind      string            `json:"kind"`
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	Labels    map[string]string `json:"labels"`
}

type Manifest struct {
	Content []byte
}

func (m *Manifest) String() string {
	return string(m.Content)
}

type InfraOutput struct {
	Values     map[string]interface{}
	HasChanges bool
}
