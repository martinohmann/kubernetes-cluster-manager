package kubernetes

import "reflect"

// ResourceSelector is used to select kubernetes resources.
type ResourceSelector struct {
	Kind      string            `json:"kind,omitempty" yaml:"kind,omitempty"`
	Name      string            `json:"name,omitempty" yaml:"name,omitempty"`
	Namespace string            `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	Labels    map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
}

// Matches returns true if other matches s.
func (s ResourceSelector) Matches(other ResourceSelector) bool {
	if s.Kind != other.Kind || s.Namespace != other.Namespace {
		return false
	}

	if s.Name != "" || other.Name != "" {
		return s.Name == other.Name
	}

	return reflect.DeepEqual(s.Labels, other.Labels)
}
