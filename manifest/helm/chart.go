package helm

import (
	"os"
	"os/exec"
	"strings"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/command"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/fs"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

// Chart defines the type for a helm chart.
type Chart struct {
	name     string
	executor command.Executor
}

// NewChart creates a new Chart value for helm chart with name. Passed in
// executor will be used to run helm commands.
func NewChart(name string, executor command.Executor) *Chart {
	return &Chart{
		name:     name,
		executor: executor,
	}
}

// Render renders the helm chart using the values from passed valueFile.
// Returns a kubernetes manifest.
func (c *Chart) Render(values map[string]interface{}) ([]byte, error) {
	content, err := yaml.Marshal(values)
	if err != nil {
		return nil, err
	}

	f, err := fs.NewTempFile("values.yaml", content)
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

	out, err := c.executor.RunSilently(cmd)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to render manifest: %s", strings.Trim(out, "\n"))
	}

	return []byte(out), nil
}
