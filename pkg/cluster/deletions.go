package cluster

import (
	"os"
	"path/filepath"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kubernetes"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/manifest"
	"github.com/pkg/errors"
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
	o *Options,
	kubectl *kubernetes.Kubectl,
	resources []ResourceSelector,
) ([]ResourceSelector, error) {
	if o.DryRun && len(resources) > 0 {
		buf, _ := yaml.Marshal(resources)
		log.Warnf("Would delete the following resources:\n%s", string(buf))
		return resources, nil
	}

	return kubectl.DeleteResources(resources)
}

func deleteManifest(o *Options, kubectl *kubernetes.Kubectl, manifest *manifest.Manifest) error {
	filename := filepath.Join(o.ManifestsDir, manifest.Filename())

	if o.DryRun {
		log.Warnf("Would delete manifest %s", filename)
		log.Debug(string(manifest.Content))
	} else {
		log.Infof("Deleting manifest %s", filename)
		if err := kubectl.DeleteManifest(manifest.Content); err != nil {
			return err
		}

		err := os.Remove(filename)
		if err != nil && !os.IsNotExist(err) {
			return errors.WithStack(err)
		}
	}

	return nil
}
