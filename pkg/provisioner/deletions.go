package provisioner

import (
	"io/ioutil"
	"os"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/api"
	"gopkg.in/yaml.v2"
)

func loadDeletions(deletionsFile string) (*api.Deletions, error) {
	deletions := api.Deletions{}

	content, err := ioutil.ReadFile(deletionsFile)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	err = yaml.Unmarshal(content, &deletions)

	return &deletions, err
}

func saveDeletions(deletionsFile string, deletions *api.Deletions) error {
	content, err := yaml.Marshal(deletions)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(deletionsFile, content, 0660)
}

func processResourceDeletions(kubectl *Kubectl, deletions []*api.Deletion) error {
	for _, deletion := range deletions {
		if err := kubectl.DeleteResource(deletion); err != nil {
			return err
		}

		deletion.MarkDeleted()
	}

	return nil
}
