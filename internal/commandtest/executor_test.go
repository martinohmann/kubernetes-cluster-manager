package commandtest

import (
	"os/exec"
	"testing"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/command"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithMockExecutor(t *testing.T) {
	def := command.DefaultExecutor
	called := false

	WithMockExecutor(func(executor MockExecutor) {
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

	WithMockExecutor(func(executor MockExecutor) {
		called = true

		panic("whoops")
	})
}

func TestWithMockExecutorWrapped(t *testing.T) {
	wrapped := NewMockExecutor(nil)

	wrapped.ExpectCommand("foo").WillReturn("bar")

	WithMockExecutor(func(executor MockExecutor) {
		executor.ExpectCommand("foo").WillExecute()

		cmd := exec.Command("foo")

		out, err := command.RunSilently(cmd)

		assert.NoError(t, err)
		assert.Equal(t, "bar", out)
	}, wrapped)
}

func TestMockExecutorExpectationsWereMet(t *testing.T) {
	e := NewMockExecutor(nil)

	e.ExpectCommand("foo").WillReturn("foo")

	e.Run(exec.Command("foo"))

	assert.NoError(t, e.ExpectationsWereMet())
}

func TestMockExecutorExpectationsWereMetError(t *testing.T) {
	e := NewMockExecutor(nil)

	e.ExpectCommand("foo").WillReturn("foo")

	assert.Error(t, e.ExpectationsWereMet())
}

func TestMockExecutorMultipleExpectations(t *testing.T) {
	e := NewMockExecutor(nil)

	e.ExpectCommand("foo").WillReturnError(errors.New("foo error"))
	e.ExpectCommand("bar").WillReturn("bar")

	_, err := e.Run(exec.Command("bar"))

	assert.NoError(t, err)
	assert.Error(t, e.ExpectationsWereMet())
}

func TestMockExecutorMultipleExpectationsAlreadyMet(t *testing.T) {
	e := NewMockExecutor(nil)

	expectedError := `all expectations were already fulfilled, command "bar" was not expected`

	e.ExpectCommand("foo").WillReturn("foo")

	_, err := e.Run(exec.Command("foo"))
	assert.NoError(t, err)

	_, err = e.Run(exec.Command("bar"))
	require.Error(t, err)
	assert.Equal(t, expectedError, err.Error())

	assert.NoError(t, e.ExpectationsWereMet())
}

func TestMockExecutorMultipleExpectationsUnordered(t *testing.T) {
	e := NewMockExecutor(nil)

	e.ExpectCommand("foo").WillReturn("foo")
	e.ExpectCommand("bar").WillReturn("bar")

	_, err := e.Run(exec.Command("bar"))
	assert.NoError(t, err)

	_, err = e.Run(exec.Command("foo"))
	assert.NoError(t, err)

	assert.NoError(t, e.ExpectationsWereMet())
}

func TestMockExecutorCommandMismatch(t *testing.T) {
	e := NewMockExecutor(nil)

	e.ExpectCommand("foo").WillReturn("foo")

	expectedError := `command "somecommand somearg" was not expected`

	_, err := e.Run(exec.Command("somecommand", "somearg"))

	assert.Equal(t, expectedError, err.Error())
}

func TestMockExecutorCommandError(t *testing.T) {
	e := NewMockExecutor(nil)

	expectedError := errors.New("foo error")

	e.ExpectCommand("foo bar").WillReturnError(expectedError)

	_, err := e.Run(exec.Command("foo", "bar"))

	assert.Equal(t, expectedError, err)
}

func TestMockExecutorCommandUnexpected(t *testing.T) {
	e := NewMockExecutor(nil)

	e.ExpectCommand("foo").WillReturn("foo")
	_, err := e.RunSilently(exec.Command("foo"))

	assert.NoError(t, err)

	_, err = e.RunSilently(exec.Command("unexpected", "command"))

	assert.Error(t, err)
}

func TestMockExecutorExecutorError(t *testing.T) {
	e := NewMockExecutor(nil)

	e.ExpectCommand("foo bar").WillExecute()

	expectedError := `cannot execute command "foo bar" because there is no executor defined`

	_, err := e.Run(exec.Command("foo", "bar"))

	assert.Equal(t, expectedError, err.Error())
}
