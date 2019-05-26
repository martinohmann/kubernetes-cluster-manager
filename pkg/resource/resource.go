package resource

import (
	"fmt"
	"strings"
)

// Resource is a kubernetes resource.
type Resource struct {
	Kind      string
	Name      string
	Namespace string
	Content   []byte
}

// Head defines the yaml structure of a resource head. This is used
// for parsing metadata from raw yaml documents.
type Head struct {
	Kind     string   `yaml:"kind"`
	Metadata Metadata `yaml:"metadata"`
}

// Metadata is the resource metadata we are interested in.
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

// String implements fmt.Stringer
func (r *Resource) String() string {
	return fmt.Sprintf("%s/%s", strings.ToLower(r.Kind), r.Name)
}

// matches returns true if other matches r. Two resources match if their name,
// kind and namespace are the same.
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
