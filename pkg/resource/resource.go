package resource

// Resource is a kubernetes resource.
type Resource struct {
	Kind      string
	Name      string
	Namespace string
	Content   []byte
}

// Head defines the yaml structure of a manifest resource head. This is used
// for parsing metadata from raw yaml documents.
type Head struct {
	Kind     string   `yaml:"kind"`
	Metadata Metadata `yaml:"metadata"`
}

type Metadata struct {
	Name        string            `yaml:"name"`
	Namespace   string            `yaml:"namespace"`
	Annotations map[string]string `yaml:"annotations"`
}

// New creates a new resource value with content and head.
func New(content []byte, head Head) *Resource {
	return &Resource{
		Kind:      head.Kind,
		Name:      head.Metadata.Name,
		Namespace: head.Metadata.Namespace,
		Content:   content,
	}
}

func (r *Resource) matches(other *Resource) bool {
	if r == other {
		return true
	}

	if r == nil || other == nil {
		return false
	}

	if r.Kind != other.Kind || r.Namespace != other.Namespace {
		return false
	}

	return r.Name == other.Name
}
