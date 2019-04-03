package terraform

import (
	"bytes"
	"encoding/json"
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

func (m *Manager) Apply() (*api.InfraOutput, error) {
	if m.cfg.DryRun {
		return m.plan()
	}

	return m.apply()
}

func (m *Manager) apply() (*api.InfraOutput, error) {
	args := terraform("apply")

	if m.cfg.Terraform.Parallelism != 0 {
		args = append(args, fmt.Sprintf("-parallelism=%d", m.cfg.Terraform.Parallelism))
	}

	if m.cfg.Terraform.AutoApprove {
		args = append(args, "-auto-approve")
	}

	err := executor.Execute(m.w, args...)
	if err != nil {
		return nil, err
	}

	return m.fetchOutputValues(&api.InfraOutput{})
}

func (m *Manager) plan() (*api.InfraOutput, error) {
	args := terraform("plan", "-detailed-exitcode")

	err := executor.Execute(m.w, args...)

	output := &api.InfraOutput{}

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 2 {
			output.HasChanges = true
			err = nil
		}
	}

	if err != nil {
		return nil, err
	}

	return m.fetchOutputValues(output)
}

func (m *Manager) fetchOutputValues(out *api.InfraOutput) (*api.InfraOutput, error) {
	var buf bytes.Buffer

	err := executor.Execute(&buf, terraform("output", "-json")...)
	if err != nil {
		return nil, err
	}

	values := make(terraformOutput)
	if err := json.Unmarshal(buf.Bytes(), &values); err != nil {
		return nil, err
	}

	out.Values = make(map[string]interface{})
	for key, ov := range values {
		out.Values[key] = ov.Value
	}

	return out, nil
}
