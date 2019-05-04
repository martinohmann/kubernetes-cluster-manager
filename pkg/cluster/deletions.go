package cluster

import (
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kubernetes"
	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

// Deletions defines the structure of a resource deletions file's content.
type Deletions struct {
	PreApply   []*kubernetes.ResourceSelector `json:"preApply" yaml:"preApply"`
	PostApply  []*kubernetes.ResourceSelector `json:"postApply" yaml:"postApply"`
	PreDestroy []*kubernetes.ResourceSelector `json:"preDestroy" yaml:"preDestroy"`
}

func processResourceDeletions(
	o *Options,
	kubectl *kubernetes.Kubectl,
	resources []*kubernetes.ResourceSelector,
) ([]*kubernetes.ResourceSelector, error) {
	if o.DryRun && len(resources) > 0 {
		buf, _ := yaml.Marshal(resources)
		log.Warnf("Would delete the following resources:\n%s", string(buf))
		return resources, nil
	}

	return kubectl.DeleteResources(resources)
}
