package commandtest

import (
	"errors"
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

func TestMockExecutorCommandMismatch(t *testing.T) {
	e := NewMockExecutor(nil)

	e.Command("foo").WillReturn("foo")

	expectedError := errors.New(`command "somecommand somearg" does not match "foo"`)

	_, err := e.RunSilently(exec.Command("somecommand", "somearg"))

	assert.Equal(t, expectedError, err)
}

func TestMockExecutorCommandPatternMismatch(t *testing.T) {
	e := NewMockExecutor(nil)

	e.Pattern("^foo$").WillReturn("foo")

	expectedError := errors.New(`command "somecommand somearg" does not match pattern "^foo$"`)

	_, err := e.RunSilently(exec.Command("somecommand", "somearg"))

	assert.Equal(t, expectedError, err)
}

func TestMockExecutorCommandUnexpected(t *testing.T) {
	e := NewMockExecutor(nil)

	e.Command("foo").WillReturn("foo")
	_, err := e.RunSilently(exec.Command("foo"))

	assert.NoError(t, err)

	_, err = e.RunSilently(exec.Command("unexpected", "command"))

	assert.Error(t, err)
}

func TestMockExecutorNextCommand(t *testing.T) {
	e := NewMockExecutor(nil)

	e.NextCommand().WillReturn("foo")

	out, err := e.RunSilently(exec.Command("somecommand"))

	assert.NoError(t, err)
	assert.Equal(t, "foo", out)

	e.NextCommand().WillSucceed()

	_, err = e.RunSilently(exec.Command("somecommand"))

	assert.NoError(t, err)
}

func TestMockExecutorAnyCommand(t *testing.T) {
	e := NewMockExecutor(nil)

	e.AnyCommand().WillReturn("foo")

	out, err := e.RunSilently(exec.Command("somecommand"))

	assert.NoError(t, err)
	assert.Equal(t, "foo", out)

	out, err = e.RunSilently(exec.Command("someothercommand"))

	assert.NoError(t, err)
	assert.Equal(t, "foo", out)
}

func TestMockExecutorPattern(t *testing.T) {
	e := NewMockExecutor(nil)

	expectedError := errors.New("some error")

	e.Pattern("^foo .*$").WillReturnError(expectedError)

	_, err := e.RunSilently(exec.Command("foo", "bar"))

	assert.Equal(t, expectedError, err)
}
