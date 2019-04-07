package api

import (
	"gopkg.in/yaml.v2"
)

// Deletions defines the structure of a resource deletions file's content.
type Deletions struct {
	PreApply   []*Deletion `json:"preApply" yaml:"preApply"`
	PostApply  []*Deletion `json:"postApply" yaml:"postApply"`
	PreDestroy []*Deletion `json:"preDestroy" yaml:"preDestroy"`
}

// FilterPending filters for all preApply, postApply and preDestroy deletions
// that are still pending and returns them.
func (d *Deletions) FilterPending() *Deletions {
	return &Deletions{
		PreApply:   filterPendingDeletions(d.PreApply),
		PostApply:  filterPendingDeletions(d.PostApply),
		PreDestroy: filterPendingDeletions(d.PreDestroy),
	}
}

// filterPendingDeletions filters for deletions that are still pending and
// returns them.
func filterPendingDeletions(deletions []*Deletion) []*Deletion {
	p := make([]*Deletion, 0)

	for _, d := range deletions {
		if !d.deleted {
			p = append(p, d)
		}
	}

	return p
}

// Deletion is the structure of a resource deletion entry.
type Deletion struct {
	Kind      string            `json:"kind" yaml:"kind"`
	Name      string            `json:"name" yaml:"name"`
	Namespace string            `json:"namespace" yaml:"namespace"`
	Labels    map[string]string `json:"labels" yaml:"labels"`

	// internal marker for sucessful resource deletion.
	deleted bool
}

// String implements fmt.Stringer.
func (d *Deletion) String() string {
	buf, _ := yaml.Marshal(d)
	return string(buf)
}

// MarkDeleted marks a deletion as successfully deleted. This indicator can be
// used to filter for deletions that are still pending.
func (d *Deletion) MarkDeleted() {
	d.deleted = true
}

// Manifest defines a Kubernetes manifest.
type Manifest struct {
	Content []byte `json:"content" yaml:"content"`
}

// String implements fmt.Stringer.
func (m *Manifest) String() string {
	return string(m.Content)
}

// InfraOutput contains the output variables of an infrastructure manager.
type InfraOutput struct {
	Values map[string]interface{} `json:"values" yaml:"values"`
}
