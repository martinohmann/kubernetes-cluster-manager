package provisioner

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/command"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	noTerraformRootModulePattern = ".*The module root could not be found. There is nothing to output.*"
)

var (
	noTerraformRootModuleRegexp = regexp.MustCompile(noTerraformRootModulePattern)
)

type terraformOutputValue struct {
	Value interface{} `json:"value"`
}

// Terraform is an infrastructure manager that uses terraform to manage
// resources.
type Terraform struct {
	Parallelism int
}

// NewTerraform creates a new terraform infrastructure manager.
func NewTerraform(o *Options) Provisioner {
	return &Terraform{
		Parallelism: o.Parallelism,
	}
}

// Provision implements Provision from the Provisioner interface.
func (m *Terraform) Provision(ctx context.Context) error {
	args := []string{
		"terraform",
		"apply",
		"--auto-approve",
	}

	if m.Parallelism > 0 {
		args = append(args, fmt.Sprintf("--parallelism=%d", m.Parallelism))
	}

	cmd := exec.Command(args[0], args[1:]...)

	_, err := command.RunWithContext(ctx, cmd)

	return err
}

// Reconcile implements Reconciler.
func (m *Terraform) Reconcile(ctx context.Context) (err error) {
	args := []string{
		"terraform",
		"plan",
		"--detailed-exitcode",
	}

	if m.Parallelism > 0 {
		args = append(args, fmt.Sprintf("--parallelism=%d", m.Parallelism))
	}

	cmd := exec.Command(args[0], args[1:]...)

	if _, err = command.RunWithContext(ctx, cmd); err != nil {
		// ExitCode 2 means that there are infrastructure changes. This is not an error.
		if exitErr, ok := errors.Cause(err).(*exec.ExitError); ok && exitErr.ExitCode() == 2 {
			err = nil
		}
	}

	return
}

// Output implements Outputter.
func (m *Terraform) Output(ctx context.Context) (map[string]interface{}, error) {
	args := []string{
		"terraform",
		"output",
		"--json",
	}

	cmd := exec.Command(args[0], args[1:]...)

	v := make(map[string]interface{})

	out, err := command.RunSilently(cmd)
	if err != nil {
		// If there was no tfstate written yet and we try to fetch output
		// variables from terraform it will fail with an error. In that case we
		// ignore the error and just return empty values.
		if noTerraformRootModuleRegexp.MatchString(err.Error()) {
			log.Warn("terraform root module was not found, most likely the tfstate was not written yet.")
			return v, nil
		}

		return nil, err
	}

	outputValues := make(map[string]terraformOutputValue)
	if err := json.Unmarshal([]byte(out), &outputValues); err != nil {
		return nil, err
	}

	for key, ov := range outputValues {
		v[key] = ov.Value
	}

	return v, nil
}

// Destroy implements Destroy from the Provisioner interface.
func (m *Terraform) Destroy(ctx context.Context) error {
	args := []string{
		"terraform",
		"destroy",
		"--auto-approve",
	}

	if m.Parallelism > 0 {
		args = append(args, fmt.Sprintf("--parallelism=%d", m.Parallelism))
	}

	cmd := exec.Command(args[0], args[1:]...)

	_, err := command.RunWithContext(ctx, cmd)

	return err
}
