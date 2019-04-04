package api

type Deletions struct {
	PreApply  []Deletion
	PostApply []Deletion
}

type Deletion struct {
	Kind      string
	Name      string
	Namespace string
	Labels    map[string]string
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
