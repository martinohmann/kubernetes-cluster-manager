package command

import (
	"errors"
	"os/exec"
	"strings"
)

type MockExecutor struct {
	willError  bool
	err        error
	willReturn bool
	out        string

	ExecutedCommands []string
}

func NewMockExecutor() *MockExecutor {
	return &MockExecutor{
		ExecutedCommands: make([]string, 0),
	}
}

func (e *MockExecutor) WillError() *MockExecutor {
	return e.WillErrorWith(errors.New("error"))
}

func (e *MockExecutor) WillErrorWith(err error) *MockExecutor {
	e.willError = true
	e.err = err

	return e
}

func (e *MockExecutor) WillReturn(out string) *MockExecutor {
	e.willReturn = true
	e.out = out

	return e
}

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

func (e *MockExecutor) RunSilently(cmd *exec.Cmd) (string, error) {
	return e.Run(cmd)
}
