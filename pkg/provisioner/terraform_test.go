package provisioner

import (
	"testing"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/command"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kcm"
	"github.com/stretchr/testify/assert"
)

func TestTerraformProvision(t *testing.T) {
	executor := command.NewMockExecutor(nil)

	options := kcm.TerraformOptions{Parallelism: 4}

	m := NewTerraform(&options, executor)

	err := m.Provision()

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
}

func TestTerraformPlan(t *testing.T) {
	executor := command.NewMockExecutor(nil)

	m := NewTerraform(&kcm.TerraformOptions{}, executor)

	err := m.Reconcile()

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
}

func TestTerraformFetch(t *testing.T) {
	executor := command.NewMockExecutor(nil)

	m := NewTerraform(&kcm.TerraformOptions{}, executor)

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

	values, err := m.Fetch()

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
}

func TestTerraformDestroy(t *testing.T) {
	executor := command.NewMockExecutor(nil)

	options := kcm.TerraformOptions{Parallelism: 4}

	m := NewTerraform(&options, executor)

	err := m.Destroy()

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
}
