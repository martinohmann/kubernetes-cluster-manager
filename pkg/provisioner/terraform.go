package provisioner

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/command"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kcm"
	"github.com/pkg/errors"
)

const (
	noTerraformRootModulePattern = ".*The module root could not be found. There is nothing to output.*"
)

var (
	noTerraformRootModuleRegexp = regexp.MustCompile(noTerraformRootModulePattern)
)

func init() {
	Register("terraform", func(o *kcm.ProvisionerOptions) (kcm.Provisioner, error) {
		return NewTerraform(&o.Terraform), nil
	})
}

type terraformOutputValue struct {
	Value interface{} `json:"value"`
}

// Terraform is an infrastructure manager that uses terraform to manage
// resources.
type Terraform struct {
	options *kcm.TerraformOptions
}

// NewTerraform creates a new terraform infrastructure manager.
func NewTerraform(o *kcm.TerraformOptions) *Terraform {
	return &Terraform{
		options: o,
	}
}

// Provision implements Provision from the kcm.Provisioner interface.
func (m *Terraform) Provision() error {
	args := []string{
		"terraform",
		"apply",
		"--auto-approve",
	}

	if m.options.Parallelism > 0 {
		args = append(args, fmt.Sprintf("--parallelism=%d", m.options.Parallelism))
	}

	cmd := exec.Command(args[0], args[1:]...)

	_, err := command.Run(cmd)

	return err
}

// Reconcile implements kcm.Reconciler.
func (m *Terraform) Reconcile() (err error) {
	args := []string{
		"terraform",
		"plan",
		"--detailed-exitcode",
	}

	if m.options.Parallelism > 0 {
		args = append(args, fmt.Sprintf("--parallelism=%d", m.options.Parallelism))
	}

	cmd := exec.Command(args[0], args[1:]...)

	if _, err = command.Run(cmd); err != nil {
		// ExitCode 2 means that there are infrastructure changes. This is not an error.
		if exitErr, ok := errors.Cause(err).(*exec.ExitError); ok && exitErr.ExitCode() == 2 {
			err = nil
		}
	}

	return
}

// Fetch implements kcm.ValueFetcher.
func (m *Terraform) Fetch() (kcm.Values, error) {
	args := []string{
		"terraform",
		"output",
		"--json",
	}

	cmd := exec.Command(args[0], args[1:]...)

	out, err := command.RunSilently(cmd)
	if err != nil {
		// If there was no tfstate written yet and we try to fetch output
		// variables from terraform it will fail with an error. In that case we
		// ignore the error and just return empty values.
		if noTerraformRootModuleRegexp.MatchString(err.Error()) {
			return kcm.Values{}, nil
		}

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

	if m.options.Parallelism > 0 {
		args = append(args, fmt.Sprintf("--parallelism=%d", m.options.Parallelism))
	}

	cmd := exec.Command(args[0], args[1:]...)

	_, err := command.Run(cmd)

	return err
}
