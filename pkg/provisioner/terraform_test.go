package provisioner

import (
	"context"
	"errors"
	"testing"

	"github.com/fatih/color"
	"github.com/martinohmann/kubernetes-cluster-manager/internal/commandtest"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kcm"
	"github.com/stretchr/testify/assert"
)

func TestTerraformProvision(t *testing.T) {
	commandtest.WithMockExecutor(func(executor *commandtest.MockExecutor) {
		options := &Options{Parallelism: 4}

		m := NewTerraform(options)

		err := m.Provision(context.Background())

		if !assert.NoError(t, err) {
			return
		}

		if assert.Len(t, executor.ExecutedCommands, 1) {
			assert.Equal(
				t,
				"terraform apply --auto-approve --parallelism=4",
				executor.ExecutedCommands[0],
			)
		}
	})
}

func TestTerraformReconcile(t *testing.T) {
	commandtest.WithMockExecutor(func(executor *commandtest.MockExecutor) {
		m := &Terraform{}

		err := m.Reconcile(context.Background())

		if !assert.NoError(t, err) {
			return
		}

		if assert.Len(t, executor.ExecutedCommands, 1) {
			assert.Equal(
				t,
				"terraform plan --detailed-exitcode",
				executor.ExecutedCommands[0],
			)
		}
	})
}

func TestTerraformOutput(t *testing.T) {
	commandtest.WithMockExecutor(func(executor *commandtest.MockExecutor) {
		m := &Terraform{}

		output := `
{
  "foo": {
	"value": "bar"
  },
  "bar": {
	"value": ["baz"]
  }
}`

		executor.NextCommand().WillReturn(output)

		expectedValues := kcm.Values{
			"foo": "bar",
			"bar": []interface{}{"baz"},
		}

		values, err := m.Output(context.Background())

		if !assert.NoError(t, err) {
			return
		}

		if assert.Len(t, executor.ExecutedCommands, 1) {
			assert.Equal(
				t,
				"terraform output --json",
				executor.ExecutedCommands[0],
			)

			assert.Equal(t, expectedValues, values)
		}
	})
}

// Ref: https://github.com/martinohmann/kubernetes-cluster-manager/issues/21
func TestTerraformOutputIssue21(t *testing.T) {
	commandtest.WithMockExecutor(func(executor *commandtest.MockExecutor) {
		m := &Terraform{}

		executor.NextCommand().WillReturnError(
			errors.New(color.RedString("The module root could not be found. There is nothing to output.")),
		)

		values, err := m.Output(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, kcm.Values{}, values)
	})
}

func TestTerraformDestroy(t *testing.T) {
	commandtest.WithMockExecutor(func(executor *commandtest.MockExecutor) {
		options := &Options{Parallelism: 4}

		m := NewTerraform(options)

		err := m.Destroy(context.Background())

		if !assert.NoError(t, err) {
			return
		}

		if assert.Len(t, executor.ExecutedCommands, 1) {
			assert.Equal(
				t,
				"terraform destroy --auto-approve --parallelism=4",
				executor.ExecutedCommands[0],
			)
		}
	})
}
