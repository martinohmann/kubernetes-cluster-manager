package terraform

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/api"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/command"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/config"
)

type outputValue struct {
	Value interface{} `json:"value"`
}

// InfraManager is an infrastructure manager that uses terraform to manage
// resources.
type InfraManager struct {
	cfg      *config.TerraformConfig
	executor command.Executor
}

// NewInfraManager creates a new terraform infrastructure manager.
func NewInfraManager(cfg *config.TerraformConfig, executor command.Executor) *InfraManager {
	return &InfraManager{
		cfg:      cfg,
		executor: executor,
	}
}

// Apply implements Apply from the api.InfraManager interface.
func (m *InfraManager) Apply() error {
	args := []string{
		"terraform",
		"apply",
		"--auto-approve",
	}

	if m.cfg.Parallelism != 0 {
		args = append(args, fmt.Sprintf("--parallelism=%d", m.cfg.Parallelism))
	}

	cmd := exec.Command(args[0], args[1:]...)

	_, err := m.executor.Run(cmd)

	return err
}

// Plan implements Plan from the api.InfraManager interface.
func (m *InfraManager) Plan() (err error) {
	args := []string{
		"terraform",
		"plan",
		"--detailed-exitcode",
	}

	if m.cfg.Parallelism != 0 {
		args = append(args, fmt.Sprintf("--parallelism=%d", m.cfg.Parallelism))
	}

	cmd := exec.Command(args[0], args[1:]...)

	if _, err = m.executor.Run(cmd); err != nil {
		// ExitCode 2 means that there are infrastructure changes. This is not an error.
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 2 {
			err = nil
		}
	}

	return
}

// GetValues implements GetValues from the api.InfraManager interface.
func (m *InfraManager) GetValues() (api.Values, error) {
	args := []string{
		"terraform",
		"output",
		"--json",
	}

	cmd := exec.Command(args[0], args[1:]...)

	out, err := m.executor.RunSilently(cmd)
	if err != nil {
		return nil, err
	}

	outputValues := make(map[string]outputValue)
	if err := json.Unmarshal([]byte(out), &outputValues); err != nil {
		return nil, err
	}

	v := make(api.Values)

	for key, ov := range outputValues {
		v[key] = ov.Value
	}

	return v, nil
}

// Destroy implements Destroy from the api.InfraManager interface.
func (m *InfraManager) Destroy() error {
	args := []string{
		"terraform",
		"destroy",
		"--auto-approve",
	}

	if m.cfg.Parallelism != 0 {
		args = append(args, fmt.Sprintf("--parallelism=%d", m.cfg.Parallelism))
	}

	cmd := exec.Command(args[0], args[1:]...)

	_, err := m.executor.Run(cmd)

	return err
}
