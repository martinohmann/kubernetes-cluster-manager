package provisioner

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/command"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kcm"
	"github.com/pkg/errors"
)

func init() {
	Register("terraform", func(o *kcm.ProvisionerOptions, e command.Executor) (kcm.Provisioner, error) {
		return NewTerraform(&o.Terraform, e), nil
	})
}

type terraformOutputValue struct {
	Value interface{} `json:"value"`
}

// Terraform is an infrastructure manager that uses terraform to manage
// resources.
type Terraform struct {
	options  *kcm.TerraformOptions
	executor command.Executor
}

// NewTerraform creates a new terraform infrastructure manager.
func NewTerraform(o *kcm.TerraformOptions, executor command.Executor) *Terraform {
	return &Terraform{
		options:  o,
		executor: executor,
	}
}

// Provision implements Provision from the kcm.Provisioner interface.
func (m *Terraform) Provision() error {
	args := []string{
		"terraform",
		"apply",
		"--auto-approve",
	}

	if m.options.Parallelism != 0 {
		args = append(args, fmt.Sprintf("--parallelism=%d", m.options.Parallelism))
	}

	cmd := exec.Command(args[0], args[1:]...)

	_, err := m.executor.Run(cmd)

	return err
}

// Reconcile implements Reconcile from the kcm.Provisioner interface.
func (m *Terraform) Reconcile() (err error) {
	args := []string{
		"terraform",
		"plan",
		"--detailed-exitcode",
	}

	if m.options.Parallelism != 0 {
		args = append(args, fmt.Sprintf("--parallelism=%d", m.options.Parallelism))
	}

	cmd := exec.Command(args[0], args[1:]...)

	if _, err = m.executor.Run(cmd); err != nil {
		// ExitCode 2 means that there are infrastructure changes. This is not an error.
		if exitErr, ok := errors.Cause(err).(*exec.ExitError); ok && exitErr.ExitCode() == 2 {
			err = nil
		}
	}

	return
}

// Fetch implements Fetch from the kcm.Provisioner interface.
func (m *Terraform) Fetch() (kcm.Values, error) {
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

	v := make(kcm.Values)

	for key, ov := range outputValues {
		v[key] = ov.Value
	}

	return v, nil
}

// Destroy implements Destroy from the kcm.Provisioner interface.
func (m *Terraform) Destroy() error {
	args := []string{
		"terraform",
		"destroy",
		"--auto-approve",
	}

	if m.options.Parallelism != 0 {
		args = append(args, fmt.Sprintf("--parallelism=%d", m.options.Parallelism))
	}

	cmd := exec.Command(args[0], args[1:]...)

	_, err := m.executor.Run(cmd)

	return err
}
