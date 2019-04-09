package terraform

import (
	"testing"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/api"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/command"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestApply(t *testing.T) {
	executor := command.NewMockExecutor()

	cfg := &config.Config{Terraform: config.TerraformConfig{Parallelism: 4}}

	m := NewInfraManager(cfg, executor)

	err := m.Apply()

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

func TestApplyDryRun(t *testing.T) {
	executor := command.NewMockExecutor()

	cfg := &config.Config{
		DryRun: true,
	}

	m := NewInfraManager(cfg, executor)

	err := m.Apply()

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

func TestGetValues(t *testing.T) {
	executor := command.NewMockExecutor()

	cfg := &config.Config{}

	m := NewInfraManager(cfg, executor)

	output := `
{
  "foo": {
	"value": "bar"
  },
  "bar": {
	"value": ["baz"]
  }
}`

	executor.WillReturn(output)

	expectedValues := api.Values{
		"foo": "bar",
		"bar": []interface{}{"baz"},
	}

	values, err := m.GetValues()

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

func TestDestroy(t *testing.T) {
	executor := command.NewMockExecutor()

	cfg := &config.Config{Terraform: config.TerraformConfig{Parallelism: 4}}

	m := NewInfraManager(cfg, executor)

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

func TestDestroyDryRun(t *testing.T) {
	executor := command.NewMockExecutor()

	cfg := &config.Config{DryRun: true}

	m := NewInfraManager(cfg, executor)

	err := m.Destroy()

	if !assert.NoError(t, err) {
		return
	}

	assert.Len(t, executor.ExecutedCommands, 0)
}
