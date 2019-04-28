package helm

import (
	"os"
	"os/exec"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/command"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/file"
	yaml "gopkg.in/yaml.v2"
)

// Chart defines the type for a helm chart.
type Chart struct {
	name string
}

// NewChart creates a new Chart value for helm chart with name.
func NewChart(name string) *Chart {
	return &Chart{
		name: name,
	}
}

// Render renders the helm chart using the values from passed valueFile.
// Returns a kubernetes manifest.
func (c *Chart) Render(values map[string]interface{}) ([]byte, error) {
	content, err := yaml.Marshal(values)
	if err != nil {
		return nil, err
	}

	f, err := file.NewTempFile("values.yaml", content)
	if err != nil {
		return nil, err
	}

	defer os.Remove(f.Name())

	args := []string{
		"helm",
		"template",
		"--values",
		f.Name(),
		c.name,
	}

	cmd := exec.Command(args[0], args[1:]...)

	out, err := command.RunSilently(cmd)
	if err != nil {
		return nil, err
	}

	return []byte(out), nil
}
