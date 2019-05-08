package commandtest

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"regexp"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/command"
)

const (
	anyCommandPattern = "^.*$"
)

var (
	anyCommandRegexp = regexp.MustCompile(anyCommandPattern)
)

// WithMockExecutor will replace the default command.Executor with a
// MockExecutor and then call f. The executor is passed to f to be able to make
// assertions and control command return values. WithMockExecutor will restore
// the previous default executor after f returns or panics. Optionally another
// executor can be passed which will be used by the MockExecutor if passthrough
// of a command is explicitly requested.
func WithMockExecutor(f func(*MockExecutor), wrapped ...command.Executor) {
	var wrappedE command.Executor
	if len(wrapped) > 0 {
		wrappedE = wrapped[0]
	}

	executor := NewMockExecutor(wrappedE)
	restoreExecutor := command.SetExecutorWithRestore(executor)
	defer restoreExecutor()

	f(executor)
}

// MockExecutor can be used in tests to stub out command execution.
type MockExecutor struct {
	executor     command.Executor
	expectation  *expectation
	expectations []*expectation
	index        int

	ExecutedCommands []string
}

// NewMockExecutor creates a new MockExecutor value.
func NewMockExecutor(executor command.Executor) *MockExecutor {
	return &MockExecutor{
		executor:         executor,
		ExecutedCommands: make([]string, 0),
	}
}

// Run implements Run from the command.Executor interface.
func (e *MockExecutor) Run(cmd *exec.Cmd) (string, error) {
	return e.RunWithContext(context.Background(), cmd)
}

// RunWithContext implements RunWithContext from the command.Executor interface.
func (e *MockExecutor) RunWithContext(ctx context.Context, cmd *exec.Cmd) (out string, err error) {
	commandLine := command.Line(cmd)

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

			out, err = e.executor.RunWithContext(ctx, cmd)
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
	commandLine := command.Line(cmd)

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

// RunSilently implements RunSilently from the command.Executor interface.
func (e *MockExecutor) RunSilently(cmd *exec.Cmd) (string, error) {
	return e.RunSilentlyWithContext(context.Background(), cmd)
}

// RunSilentlyWithContext implements RunSilentlyWithContext from the command.Executor interface.
func (e *MockExecutor) RunSilentlyWithContext(ctx context.Context, cmd *exec.Cmd) (string, error) {
	return e.RunWithContext(ctx, cmd)
}

// AnyCommand creates an expectation for any command that will be executed.
func (e *MockExecutor) AnyCommand() *expectation {
	ex := &expectation{pattern: anyCommandPattern, re: anyCommandRegexp}
	e.expectation = ex
	e.expectations = nil
	e.index = 0

	return ex
}

// NextCommand creates an expectation for the next command to be executed.
func (e *MockExecutor) NextCommand() *expectation {
	ex := &expectation{pattern: anyCommandPattern, re: anyCommandRegexp}
	e.expectation = nil
	e.addExpectation(ex)

	return ex
}

// Pattern creates an expectation for an exact cmd.
func (e *MockExecutor) Command(cmd string) *expectation {
	ex := &expectation{cmd: cmd}
	e.expectation = nil
	e.addExpectation(ex)

	return ex
}

// Pattern creates an expectation for a cmd pattern.
func (e *MockExecutor) Pattern(pattern string) *expectation {
	re := regexp.MustCompile(pattern)
	ex := &expectation{pattern: pattern, re: re}
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
	cmd     string
	pattern string
	re      *regexp.Regexp
	execute bool
	err     error
	out     string
}

func (e *expectation) WillReturnError(err error) {
	e.err = err
}

func (e *expectation) WillError() {
	e.WillReturnError(errors.New("error"))
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
