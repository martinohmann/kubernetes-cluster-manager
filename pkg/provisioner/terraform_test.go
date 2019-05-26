package provisioner

import (
	"context"
	"errors"
	"testing"

	"github.com/fatih/color"
	"github.com/martinohmann/kubernetes-cluster-manager/internal/commandtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTerraformProvision(t *testing.T) {
	commandtest.WithMockExecutor(func(executor commandtest.MockExecutor) {
		options := &Options{Parallelism: 4}

		m := NewTerraform(options)

		executor.ExpectCommand("terraform apply --auto-approve --parallelism=4")

		assert.NoError(t, m.Provision(context.Background()))
		assert.NoError(t, executor.ExpectationsWereMet())
	})
}

func TestTerraformReconcile(t *testing.T) {
	commandtest.WithMockExecutor(func(executor commandtest.MockExecutor) {
		m := &Terraform{}

		executor.ExpectCommand("terraform plan --detailed-exitcode")

		assert.NoError(t, m.Reconcile(context.Background()))
		assert.NoError(t, executor.ExpectationsWereMet())
	})
}

func TestTerraformOutput(t *testing.T) {
	commandtest.WithMockExecutor(func(executor commandtest.MockExecutor) {
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

		executor.ExpectCommand("terraform output --json").WillReturn(output)

		expectedValues := map[string]interface{}{
			"foo": "bar",
			"bar": []interface{}{"baz"},
		}

		values, err := m.Output(context.Background())

		require.NoError(t, err)
		assert.Equal(t, expectedValues, values)
		assert.NoError(t, executor.ExpectationsWereMet())
	})
}

// Ref: https://github.com/martinohmann/kubernetes-cluster-manager/issues/21
func TestTerraformOutputIssue21(t *testing.T) {
	commandtest.WithMockExecutor(func(executor commandtest.MockExecutor) {
		m := &Terraform{}

		executor.ExpectCommand("terraform output --json").WillReturnError(
			errors.New(color.RedString("The module root could not be found. There is nothing to output.")),
		)

		values, err := m.Output(context.Background())

		require.NoError(t, err)
		assert.Equal(t, map[string]interface{}{}, values)
		assert.NoError(t, executor.ExpectationsWereMet())
	})
}

func TestTerraformDestroy(t *testing.T) {
	commandtest.WithMockExecutor(func(executor commandtest.MockExecutor) {
		options := &Options{Parallelism: 4}

		m := NewTerraform(options)

		executor.ExpectCommand("terraform destroy --auto-approve --parallelism=4")

		assert.NoError(t, m.Destroy(context.Background()))
		assert.NoError(t, executor.ExpectationsWereMet())
	})
}
