package provisioner

import (
	"io/ioutil"

	"github.com/martinohmann/cluster-manager/pkg/api"
	"gopkg.in/yaml.v2"
)

func loadDeletions(deletionsFile string) (*api.Deletions, error) {
	content, err := ioutil.ReadFile(deletionsFile)
	if err != nil {
		return nil, err
	}

	deletions := api.Deletions{}

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
