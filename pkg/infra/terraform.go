package infra

import (
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/api"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/command"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/config"
	log "github.com/sirupsen/logrus"
)

type terraformOutput map[string]terraformOutputValue

type terraformOutputValue struct {
	Value interface{} `json:"value"`
}

// TerraformManager is an infrastructure manager that uses terraform to manager
// resources.
type TerraformManager struct {
	cfg      *config.Config
	executor command.Executor
}

// NewTerraformManager creates a new TerraformManager value.
func NewTerraformManager(cfg *config.Config, executor command.Executor) *TerraformManager {
	return &TerraformManager{
		cfg:      cfg,
		executor: executor,
	}
}

func (m *TerraformManager) Apply() error {
	if m.cfg.DryRun {
		return m.plan()
	}

	return m.apply()
}

func (m *TerraformManager) apply() error {
	args := []string{
		"terraform",
		"apply",
		"--auto-approve",
	}

	if m.cfg.Terraform.Parallelism != 0 {
		args = append(args, fmt.Sprintf("--parallelism=%d", m.cfg.Terraform.Parallelism))
	}

	cmd := exec.Command(args[0], args[1:]...)

	_, err := m.executor.Run(cmd)

	return err
}

func (m *TerraformManager) plan() (err error) {
	args := []string{
		"terraform",
		"plan",
		"--detailed-exitcode",
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

func (m *TerraformManager) GetOutput() (*api.InfraOutput, error) {
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

	values := make(terraformOutput)
	if err := json.Unmarshal([]byte(out), &values); err != nil {
		return nil, err
	}

	output := &api.InfraOutput{}

	output.Values = make(map[string]interface{})
	for key, ov := range values {
		output.Values[key] = ov.Value
	}

	return output, nil
}

func (m *TerraformManager) Destroy() error {
	if m.cfg.DryRun {
		log.Warn("Would destroy infrastructure")
		return nil
	}

	return errors.New("destroy not implemented yet")
}
