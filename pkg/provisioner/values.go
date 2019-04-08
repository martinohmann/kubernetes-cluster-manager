package provisioner

import (
	"io/ioutil"
	"os"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/api"
	"gopkg.in/yaml.v2"
)

func loadValues(filename string) (api.Values, error) {
	v := make(api.Values)

	content, err := ioutil.ReadFile(filename)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	err = yaml.Unmarshal(content, &v)

	return v, err
}
