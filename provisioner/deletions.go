package provisioner

import (
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/api"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kubernetes"
	log "github.com/sirupsen/logrus"
)

func processResourceDeletions(o *Options, l *log.Logger, kubectl *kubernetes.Kubectl, deletions []*api.Deletion) error {
	for _, deletion := range deletions {
		if o.DryRun {
			l.Warnf("Would delete the following resource:\n%s", deletion)
			continue
		}

		if err := kubectl.DeleteResource(deletion); err != nil {
			return err
		}

		deletion.MarkDeleted()
	}

	return nil
}
