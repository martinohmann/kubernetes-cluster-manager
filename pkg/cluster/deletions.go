package cluster

import (
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kcm"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kubernetes"
	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

func processResourceDeletions(o *kcm.Options, l *log.Logger, kubectl *kubernetes.Kubectl, deletions []*kcm.Deletion) error {
	if o.DryRun && len(deletions) > 0 {
		buf, _ := yaml.Marshal(deletions)
		l.Warnf("Would delete the following resources:\n%s", string(buf))
		return nil
	}

	for _, deletion := range deletions {
		if err := kubectl.DeleteResource(deletion); err != nil {
			return err
		}

		deletion.MarkDeleted()
	}

	return nil
}
