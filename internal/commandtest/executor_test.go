package commandtest

import (
	"os/exec"
	"testing"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/command"
	"github.com/stretchr/testify/assert"
)

func TestWithMockExecutor(t *testing.T) {
	def := command.DefaultExecutor
	called := false

	WithMockExecutor(func(executor *MockExecutor) {
		called = true

		assert.Equal(t, executor, command.DefaultExecutor)
	})

	assert.True(t, called)
	assert.Equal(t, def, command.DefaultExecutor)
}

func TestWithMockExecutorPanic(t *testing.T) {
	def := command.DefaultExecutor
	called := false

	defer func() {
		recover()

		assert.True(t, called)
		assert.Equal(t, def, command.DefaultExecutor)
	}()

	WithMockExecutor(func(executor *MockExecutor) {
		called = true

		panic("whoops")
	})
}

func TestWithMockExecutorWrapped(t *testing.T) {
	wrapped := NewMockExecutor(nil)

	wrapped.Command("foo").WillReturn("bar")

	WithMockExecutor(func(executor *MockExecutor) {
		executor.Command("foo").WillExecute()

		cmd := exec.Command("foo")

		out, err := command.RunSilently(cmd)

		assert.NoError(t, err)
		assert.Equal(t, "bar", out)
	}, wrapped)
}
