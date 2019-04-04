package terraform

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/martinohmann/cluster-manager/pkg/api"
	"github.com/martinohmann/cluster-manager/pkg/config"
	"github.com/martinohmann/cluster-manager/pkg/executor"
	"github.com/martinohmann/cluster-manager/pkg/infra"
)

var _ infra.Manager = &Manager{}

type terraformOutput map[string]outputValue

type outputValue struct {
	Type      string
	Sensitive bool
	Value     interface{}
}

type Manager struct {
	cfg *config.Config
	w   io.Writer
}

func NewInfraManager(cfg *config.Config, w io.Writer) *Manager {
	if w == nil {
		w = os.Stdout
	}

	return &Manager{
		w:   w,
		cfg: cfg,
	}
}

func (m *Manager) Apply() error {
	if m.cfg.DryRun {
		return m.plan()
	}

	return m.apply()
}

func (m *Manager) apply() error {
	args := []string{
		"terraform",
		"apply",
	}

	if m.cfg.Terraform.Parallelism != 0 {
		args = append(args, fmt.Sprintf("-parallelism=%d", m.cfg.Terraform.Parallelism))
	}

	if m.cfg.Terraform.AutoApprove {
		args = append(args, "-auto-approve")
	}

	return executor.Execute(m.w, args...)
}

func (m *Manager) plan() error {
	args := []string{
		"terraform",
		"plan",
		"-detailed-exitcode",
	}

	err := executor.Execute(m.w, args...)

	if err != nil {
		// ExitCode 2 means that there are infrastructure changes. This is not an error.
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 2 {
			err = nil
		}
	}

	return err
}

func (m *Manager) GetOutput() (*api.InfraOutput, error) {
	var buf bytes.Buffer

	args := []string{
		"terraform",
		"output",
		"-json",
	}

	err := executor.Execute(&buf, args...)
	if err != nil {
		return nil, err
	}

	values := make(terraformOutput)
	if err := json.Unmarshal(buf.Bytes(), &values); err != nil {
		return nil, err
	}

	output := &api.InfraOutput{}

	output.Values = make(map[string]interface{})
	for key, ov := range values {
		output.Values[key] = ov.Value
	}

	return output, nil
}

func (m *Manager) Destroy() error {
	return errors.New("destroy not implemented yet")
}
