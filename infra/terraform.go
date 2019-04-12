package infra

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/api"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/command"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/config"
)

type terraformOutputValue struct {
	Value interface{} `json:"value"`
}

// TerraformManager is an infrastructure manager that uses terraform to manage
// resources.
type TerraformManager struct {
	cfg      *config.TerraformConfig
	executor command.Executor
}

// NewTerraformManager creates a new terraform infrastructure manager.
func NewTerraformManager(cfg *config.TerraformConfig, executor command.Executor) *TerraformManager {
	return &TerraformManager{
		cfg:      cfg,
		executor: executor,
	}
}

// Apply implements Apply from the Manager interface.
func (m *TerraformManager) Apply() error {
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

// Plan implements Plan from the Manager interface.
func (m *TerraformManager) Plan() (err error) {
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

// GetValues implements GetValues from the Manager interface.
func (m *TerraformManager) GetValues() (api.Values, error) {
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

	outputValues := make(map[string]terraformOutputValue)
	if err := json.Unmarshal([]byte(out), &outputValues); err != nil {
		return nil, err
	}

	v := make(api.Values)

	for key, ov := range outputValues {
		v[key] = ov.Value
	}

	return v, nil
}

// Destroy implements Destroy from the Manager interface.
func (m *TerraformManager) Destroy() error {
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
