package helm

import (
	"os/exec"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/api"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/command"
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
func (c *Chart) Render(valuesFile string) (*api.Manifest, error) {
	args := []string{
		"helm",
		"template",
		"--values",
		valuesFile,
		c.name,
	}

	cmd := exec.Command(args[0], args[1:]...)

	out, err := c.executor.RunSilently(cmd)
	if err != nil {
		return nil, err
	}

	return api.NewManifestFromString(out), nil
}
