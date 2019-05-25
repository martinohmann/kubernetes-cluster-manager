package cluster

import (
	"context"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kubernetes"
	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

// ResourceSelector is a type alias for kubernetes.ResourceSelector.
type ResourceSelector = kubernetes.ResourceSelector

// Deletions defines the structure of a resource deletions file's content.
type Deletions struct {
	PreApply   []ResourceSelector `json:"preApply" yaml:"preApply"`
	PostApply  []ResourceSelector `json:"postApply" yaml:"postApply"`
	PreDestroy []ResourceSelector `json:"preDestroy" yaml:"preDestroy"`
}

func processResourceDeletions(
	ctx context.Context,
	o *Options,
	kubectl *kubernetes.Kubectl,
	resources []ResourceSelector,
) ([]ResourceSelector, error) {
	if o.DryRun && len(resources) > 0 {
		buf, _ := yaml.Marshal(resources)
		log.Warnf("Would delete the following resources:\n%s", string(buf))
		return resources, nil
	}

	return kubectl.DeleteResources(ctx, resources)
}
