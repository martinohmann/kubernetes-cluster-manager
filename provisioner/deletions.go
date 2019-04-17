package provisioner

import (
	"io/ioutil"
	"os"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/api"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kubernetes"
	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

func loadDeletions(filename string) (*api.Deletions, error) {
	deletions := api.Deletions{}

	content, err := ioutil.ReadFile(filename)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	err = yaml.Unmarshal(content, &deletions)

	return &deletions, err
}

func processResourceDeletions(o *Options, kubectl *kubernetes.Kubectl, deletions []*api.Deletion) error {
	for _, deletion := range deletions {
		if o.DryRun {
			log.Warnf("Would delete the following resource:\n%s", deletion)
			continue
		}

		if err := kubectl.DeleteResource(deletion); err != nil {
			return err
		}

		deletion.MarkDeleted()
	}

	return nil
}
