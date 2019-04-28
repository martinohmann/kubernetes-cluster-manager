package kcm

import (
	"github.com/imdario/mergo"
	yaml "gopkg.in/yaml.v2"
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
	return filterDeletions(deletions, func(d *Deletion) bool { return !d.deleted })
}

// filterDeletions filters deletions using a filter func.
func filterDeletions(deletions []*Deletion, f func(*Deletion) bool) []*Deletion {
	p := make([]*Deletion, 0)

	for _, d := range deletions {
		if f(d) {
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

	// internal marker for successful resource deletion.
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

// Deleted returns true if the resource deletion executed successfully.
func (d *Deletion) Deleted() bool {
	return d.deleted
}

// Manifest defines a Kubernetes manifest.
type Manifest []byte

// Values contains the output values of an infrastructure manager.
type Values map[string]interface{}

// Merge deep merges other on top of v and overrides values already present in
// v.
func (v Values) Merge(other Values) error {
	return mergo.Merge(&v, other, mergo.WithOverride)
}

// Credentials holds the credentials needed to communicate with a Kubernetes
// cluster.
type Credentials struct {
	Server     string `json:"server" yaml:"server"`
	Token      string `json:"token" yaml:"token"`
	Kubeconfig string `json:"kubeconfig" yaml:"kubeconfig"`
	Context    string `json:"context" yaml:"context"`
}