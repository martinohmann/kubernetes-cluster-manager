package provisioner

import (
	"io/ioutil"
	"os"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/api"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kubernetes"
	"gopkg.in/yaml.v2"
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

func processResourceDeletions(kubectl *kubernetes.Kubectl, deletions []*api.Deletion) error {
	for _, deletion := range deletions {
		if err := kubectl.DeleteResource(deletion); err != nil {
			return err
		}

		deletion.MarkDeleted()
	}

	return nil
}
