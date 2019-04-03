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

type InfraOutput struct {
	Values     map[string]interface{}
	HasChanges bool
}
