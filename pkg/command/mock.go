package command

import (
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

// MockExecutor can be used in tests to stub out command execution.
type MockExecutor struct {
	willError  bool
	err        error
	willReturn bool
	out        string

	expectations []*expectation
	index        int

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
	if e.expectations != nil {
		if len(e.expectations) <= e.index {
			return "", fmt.Errorf("unexpected command %v", cmd)
		}

		expectation := e.expectations[e.index]
		commandLine := strings.Join(cmd.Args, " ")

		if expectation.re != nil {
			if !expectation.re.MatchString(commandLine) {
				return "", fmt.Errorf(
					"command %q does not match pattern %q",
					commandLine,
					expectation.pattern,
				)
			}
		} else if expectation.cmd != commandLine {
			return "", fmt.Errorf(
				"command %q does not match %q",
				commandLine,
				expectation.cmd,
			)
		}

		if expectation.err != nil {
			err = expectation.err
		}

		if expectation.out != "" {
			out = expectation.out
		}

		e.index++
	}

	e.ExecutedCommands = append(e.ExecutedCommands, strings.Join(cmd.Args, " "))

	return
}

// RunSilently implements RunSilently from Executor interface.
func (e *MockExecutor) RunSilently(cmd *exec.Cmd) (string, error) {
	return e.Run(cmd)
}

func (e *MockExecutor) Command(cmd string) *expectation {
	ex := &expectation{executor: e, cmd: cmd}
	e.addExpectation(ex)

	return ex
}

func (e *MockExecutor) Pattern(pattern string) *expectation {
	re := regexp.MustCompile(pattern)
	ex := &expectation{executor: e, pattern: pattern, re: re}
	e.addExpectation(ex)

	return ex
}

func (e *MockExecutor) addExpectation(ex *expectation) {
	if e.expectations == nil {
		e.expectations = make([]*expectation, 0)
	}

	e.expectations = append(e.expectations, ex)
}

type expectation struct {
	executor *MockExecutor
	cmd      string
	pattern  string
	re       *regexp.Regexp
	err      error
	out      string
}

func (e *expectation) WillReturnError(err error) {
	e.err = err
}

func (e *expectation) WillReturn(out string) {
	e.out = out
}
