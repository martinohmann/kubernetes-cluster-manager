package api

import (
	"gopkg.in/yaml.v2"
)

type Deletions struct {
	PreApply   []Deletion `json:"preApply" yaml:"preApply"`
	PostApply  []Deletion `json:"postApply" yaml:"postApply"`
	PreDestroy []Deletion `json:"preDestroy" yaml:"preDestroy"`
}

type Deletion struct {
	Kind      string            `json:"kind" yaml:"kind"`
	Name      string            `json:"name" yaml:"name"`
	Namespace string            `json:"namespace" yaml:"namespace"`
	Labels    map[string]string `json:"labels" yaml:"labels"`
}

// String implements fmt.Stringer.
func (d Deletion) String() string {
	buf, _ := yaml.Marshal(d)
	return string(buf)
}

type Manifest struct {
	Content []byte `json:"content" yaml:"content"`
}

// String implements fmt.Stringer.
func (m *Manifest) String() string {
	return string(m.Content)
}

type InfraOutput struct {
	Values map[string]interface{} `json:"values" yaml:"values"`
}
