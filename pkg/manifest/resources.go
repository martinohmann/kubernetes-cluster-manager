package manifest

import (
	"bytes"
	"io"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kubernetes"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// ResourceSelector is a type alias for kubernetes.ResourceSelector.
type ResourceSelector = kubernetes.ResourceSelector

// resource defines the parts of a Kubernetes resource we are interested in
// when decoding a manifest.
type resource struct {
	Kind     string `yaml:"kind"`
	Metadata struct {
		Name      string `yaml:"name"`
		Namespace string `yaml:"namespace"`
	} `yaml:"metadata"`
}

// Resources returns a slice of ResourceSelector for all resources defined in
// the manifest.
func (m *Manifest) Resources() []ResourceSelector {
	resources := make([]ResourceSelector, 0)

	if m == nil {
		return resources
	}

	buf := bytes.NewBuffer(m.Content)
	d := yaml.NewDecoder(buf)

	for {
		var r resource
		err := d.Decode(&r)
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Debugf("error while decoding manifest resources: %s", err.Error())
			continue
		}

		if r.Kind == "" || r.Metadata.Name == "" {
			continue
		}

		resources = append(resources, ResourceSelector{
			Kind:      r.Kind,
			Name:      r.Metadata.Name,
			Namespace: r.Metadata.Namespace,
		})
	}

	return resources
}

// GetVanishedResources returns selectors for all resources that are not present
// in the next revision of the manifests.
func (r Revision) GetVanishedResources() []ResourceSelector {
	vanished := make([]ResourceSelector, 0)
	nextResources := r.Next.Resources()

	for _, r := range r.Prev.Resources() {
		if i := findMatchingResource(nextResources, r); i == -1 {
			vanished = append(vanished, r)
		}
	}

	return vanished
}

func findMatchingResource(haystack []ResourceSelector, needle ResourceSelector) int {
	for i, r := range haystack {
		if r.Matches(needle) {
			return i
		}
	}

	return -1
}
