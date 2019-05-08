package commandtest

import (
	"context"
	"os/exec"
	"regexp"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/command"
	"github.com/pkg/errors"
)

// WithMockExecutor will replace the default command.Executor with a
// MockExecutor and then call f. The executor is passed to f to be able to make
// assertions and control command return values. WithMockExecutor will restore
// the previous default executor after f returns or panics. Optionally another
// executor can be passed which will be used by the MockExecutor if passthrough
// of a command is explicitly requested.
func WithMockExecutor(f func(MockExecutor), wrapped ...command.Executor) {
	var wrappedE command.Executor
	if len(wrapped) > 0 {
		wrappedE = wrapped[0]
	}

	executor := NewMockExecutor(wrappedE)
	restoreExecutor := command.SetExecutorWithRestore(executor)
	defer restoreExecutor()

	f(executor)
}

// MockExecutor is a command.Executor that allows to mock executed commands.
type MockExecutor interface {
	command.Executor

	// ExpectedCommand creates a new expectation for command cmd. Cmd can be a
	// regular expression.
	ExpectCommand(cmd string) *ExpectedCommand

	// ExpectationsWereMet should be called after all test assertions have been
	// made. It returns an error if there are commands that where expected but
	// have not been called.
	ExpectationsWereMet() error
}

// mockExecutor can be used in tests to stub out command execution.
type mockExecutor struct {
	executor command.Executor
	expected []*ExpectedCommand
}

// NewMockExecutor creates a new MockExecutor.
func NewMockExecutor(executor command.Executor) MockExecutor {
	return &mockExecutor{
		executor: executor,
		expected: make([]*ExpectedCommand, 0),
	}
}

// Run implements Run from the command.Executor interface.
func (e *mockExecutor) Run(cmd *exec.Cmd) (string, error) {
	return e.RunWithContext(context.Background(), cmd)
}

// RunSilently implements RunSilently from the command.Executor interface.
func (e *mockExecutor) RunSilently(cmd *exec.Cmd) (string, error) {
	return e.RunSilentlyWithContext(context.Background(), cmd)
}

// RunWithContext implements RunWithContext from the command.Executor interface.
func (e *mockExecutor) RunWithContext(ctx context.Context, cmd *exec.Cmd) (string, error) {
	return e.run(ctx, cmd)
}

// RunSilentlyWithContext implements RunSilentlyWithContext from the command.Executor interface.
func (e *mockExecutor) RunSilentlyWithContext(ctx context.Context, cmd *exec.Cmd) (string, error) {
	return e.RunWithContext(ctx, cmd)
}

// ExpectCommand implements ExpectCommand from the MockExecutor interface.
func (e *mockExecutor) ExpectCommand(cmd string) *ExpectedCommand {
	expected := &ExpectedCommand{command: cmd}
	e.expected = append(e.expected, expected)

	return expected
}

// ExpectationsWereMet implements ExpectationsWereMet from the MockExecutor interface.
func (e *mockExecutor) ExpectationsWereMet() error {
	for _, expectation := range e.expected {
		if !expectation.fulfilled {
			return errors.Errorf("there is a remaining expectation which was not matched:\n%s", expectation)
		}
	}

	return nil
}

func (e *mockExecutor) run(ctx context.Context, cmd *exec.Cmd) (string, error) {
	commandLine := command.Line(cmd)

	ex, err := e.findExpectation(commandLine)
	if err != nil {
		return "", err
	}

	if ex.execute {
		if e.executor == nil {
			return "", errors.Errorf(
				"cannot execute command %q because there is no executor defined",
				commandLine,
			)
		}

		return e.executor.RunWithContext(ctx, cmd)
	}

	return ex.out, ex.err
}

func (e *mockExecutor) findExpectation(commandLine string) (*ExpectedCommand, error) {
	var expected *ExpectedCommand
	var fulfilled int

	for _, next := range e.expected {
		if next.fulfilled {
			fulfilled++
			continue
		}

		if err := matchCommand(next.command, commandLine); err == nil {
			expected = next
			break
		}
	}

	if expected == nil {
		msg := "command %q was not expected"
		if fulfilled == len(e.expected) {
			msg = "all expectations were already fulfilled, " + msg
		}

		return nil, errors.Errorf(msg, commandLine)
	}

	expected.fulfilled = true

	return expected, nil
}

func matchCommand(expected, actual string) error {
	re := regexp.MustCompile(expected)

	if !re.MatchString(actual) {
		return errors.Errorf("command %q does not match %q", actual, expected)
	}

	return nil
}
