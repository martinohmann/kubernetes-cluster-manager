package command

import (
	"errors"
	"fmt"
	"os/exec"
	"regexp"
)

const (
	anyCommandPattern = "^.*$"
)

var (
	anyCommandRegexp = regexp.MustCompile(anyCommandPattern)
)

// MockExecutor can be used in tests to stub out command execution.
type MockExecutor struct {
	executor     Executor
	expectation  *expectation
	expectations []*expectation
	index        int

	ExecutedCommands []string
}

// NewMockExecutor creates a new MockExecutor value.
func NewMockExecutor(executor Executor) *MockExecutor {
	return &MockExecutor{
		executor:         executor,
		ExecutedCommands: make([]string, 0),
	}
}

// Run implements Run from Executor interface.
func (e *MockExecutor) Run(cmd *exec.Cmd) (out string, err error) {
	commandLine := commandLine(cmd)

	var ex *expectation
	if e.expectation != nil {
		ex = e.expectation
	} else if e.expectations != nil {
		if len(e.expectations) <= e.index {
			return "", fmt.Errorf("unexpected command %q", commandLine)
		}

		ex = e.expectations[e.index]
		e.index++
	}

	if ex != nil {
		if err := validateExpectation(ex, cmd); err != nil {
			return "", err
		}

		if ex.execute {
			if e.executor == nil {
				return "", fmt.Errorf(
					"cannot execute command %q because there is no executor defined",
					commandLine,
				)
			}

			out, err = e.executor.Run(cmd)
		} else {
			if ex.err != nil {
				err = ex.err
			}

			if ex.out != "" {
				out = ex.out
			}
		}
	}

	e.ExecutedCommands = append(e.ExecutedCommands, commandLine)

	return
}

func validateExpectation(ex *expectation, cmd *exec.Cmd) error {
	commandLine := commandLine(cmd)

	if ex.re != nil {
		if !ex.re.MatchString(commandLine) {
			return fmt.Errorf(
				"command %q does not match pattern %q",
				commandLine,
				ex.pattern,
			)
		}
	} else if ex.cmd != commandLine {
		return fmt.Errorf(
			"command %q does not match %q",
			commandLine,
			ex.cmd,
		)
	}

	return nil
}

// RunSilently implements RunSilently from Executor interface.
func (e *MockExecutor) RunSilently(cmd *exec.Cmd) (string, error) {
	return e.Run(cmd)
}

func (e *MockExecutor) AnyCommand() *expectation {
	ex := &expectation{executor: e, pattern: anyCommandPattern, re: anyCommandRegexp}
	e.expectation = ex
	e.expectations = nil
	e.index = 0

	return ex
}

func (e *MockExecutor) NextCommand() *expectation {
	ex := &expectation{executor: e, pattern: anyCommandPattern, re: anyCommandRegexp}
	e.expectation = nil
	e.addExpectation(ex)

	return ex
}

func (e *MockExecutor) Command(cmd string) *expectation {
	ex := &expectation{executor: e, cmd: cmd}
	e.expectation = nil
	e.addExpectation(ex)

	return ex
}

func (e *MockExecutor) Pattern(pattern string) *expectation {
	re := regexp.MustCompile(pattern)
	ex := &expectation{executor: e, pattern: pattern, re: re}
	e.expectation = nil
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
	execute  bool
	err      error
	out      string
}

func (e *expectation) WillReturnError(err error) {
	e.err = err
}

func (e *expectation) WillError() {
	e.err = errors.New("error")
}

func (e *expectation) WillReturn(out string) {
	e.out = out
}

func (e *expectation) WillSucceed() {
	e.err = nil
}

func (e *expectation) WillExecute() {
	e.execute = true
}
