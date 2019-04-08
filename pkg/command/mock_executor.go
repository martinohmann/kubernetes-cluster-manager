package command

import (
	"errors"
	"os/exec"
	"strings"
)

// MockExecutor can be used in tests to stub out command execution.
type MockExecutor struct {
	willError  bool
	err        error
	willReturn bool
	out        string

	ExecutedCommands []string
}

// NewMockExecutor creates a new MockExecutor value.
func NewMockExecutor() *MockExecutor {
	return &MockExecutor{
		ExecutedCommands: make([]string, 0),
	}
}

// WillError will make the executor return an error upon next invocation.
func (e *MockExecutor) WillError() *MockExecutor {
	return e.WillErrorWith(errors.New("error"))
}

// WillErrorWith will make the executor return the provided error upon next
// invocation.
func (e *MockExecutor) WillErrorWith(err error) *MockExecutor {
	e.willError = true
	e.err = err

	return e
}

// WillReturn will make the executor return the provided output upon next
// invocation.
func (e *MockExecutor) WillReturn(out string) *MockExecutor {
	e.willReturn = true
	e.out = out

	return e
}

// Run implements Run from Executor interface.
func (e *MockExecutor) Run(cmd *exec.Cmd) (out string, err error) {
	if e.willReturn {
		e.willReturn = false
		out = e.out
	}

	if e.willError {
		e.willError = false
		err = e.err
	}

	e.ExecutedCommands = append(e.ExecutedCommands, strings.Join(cmd.Args, " "))

	return
}

// RunSilently implements RunSilently from Executor interface.
func (e *MockExecutor) RunSilently(cmd *exec.Cmd) (string, error) {
	return e.Run(cmd)
}
