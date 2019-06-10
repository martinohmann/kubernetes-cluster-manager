package command

import (
	"bytes"
	"context"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/log"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type contextKey string

// CancelSignal is the context key that can be set to indicate which signal
// should be emitted to the running process when the context is cancelled.
//
//     context.WithValue(ctx, command.CancelSignal, os.Kill)
const CancelSignal contextKey = "cancelSignal"

// Executor defines the interface for a command executor.
type Executor interface {
	// Run executes given command and returns its output.
	Run(*exec.Cmd) (string, error)

	// Run executes given command and returns its output. The context can be
	// used to send signals to the running process.
	RunWithContext(context.Context, *exec.Cmd) (string, error)

	// RunSilently executes the given command and returns its output. Will not
	// write command output to stdout or stderr.
	RunSilently(*exec.Cmd) (string, error)

	// RunSilently executes the given command and returns its output. Will not
	// write command output to stdout or stderr. The context can be used to
	// send signals to the running process.
	RunSilentlyWithContext(context.Context, *exec.Cmd) (string, error)
}

// DefaultExecutor is the default executor used in the package level Run and
// RunSilently funcs.
var DefaultExecutor = NewExecutor(nil)

// Run runs a command using the default executor.
func Run(cmd *exec.Cmd) (string, error) {
	return DefaultExecutor.Run(cmd)
}

// RunWithContext runs a command with ctx using the default executor.
func RunWithContext(ctx context.Context, cmd *exec.Cmd) (string, error) {
	return DefaultExecutor.RunWithContext(ctx, cmd)
}

// RunSilently runs cmd silently using the default executor.
func RunSilently(cmd *exec.Cmd) (string, error) {
	return DefaultExecutor.RunSilently(cmd)
}

// RunSilentlyWithContext runs cmd with ctx silently using the default executor.
func RunSilentlyWithContext(ctx context.Context, cmd *exec.Cmd) (string, error) {
	return DefaultExecutor.RunSilentlyWithContext(ctx, cmd)
}

type executor struct {
	logger *logrus.Logger
}

// NewExecutor creates a new command executor. Accepts a logger for logging
// command output. If nil is provided logrus.StandardLogger() will be used.
func NewExecutor(l *logrus.Logger) Executor {
	if l == nil {
		l = logrus.StandardLogger()
	}

	return &executor{logger: l}
}

// Run implements Run from Executor interface.
func (e *executor) Run(cmd *exec.Cmd) (string, error) {
	return e.RunWithContext(context.Background(), cmd)
}

// RunWithContext implements RunWithContext from Executor interface.
func (e *executor) RunWithContext(ctx context.Context, cmd *exec.Cmd) (string, error) {
	var out bytes.Buffer

	prefix := color.BlueString(cmd.Args[0])
	logger := e.logger.WithContext(log.ContextWithPrefix(prefix))

	cmd.Stdout = io.MultiWriter(&out, log.LineWriter(logger.Info))
	cmd.Stderr = io.MultiWriter(&out, log.LineWriter(logger.Error))

	return e.run(ctx, &out, cmd)
}

// RunSilently implements RunSilently from Executor interface.
func (e *executor) RunSilently(cmd *exec.Cmd) (out string, err error) {
	return e.RunSilentlyWithContext(context.Background(), cmd)
}

// RunSilentlyWithContext implements RunSilentlyWithContext from Executor interface.
func (e *executor) RunSilentlyWithContext(ctx context.Context, cmd *exec.Cmd) (out string, err error) {
	var buf bytes.Buffer

	cmd.Stdout = &buf
	cmd.Stderr = &buf

	out, err = e.run(ctx, &buf, cmd)
	if err != nil {
		err = errors.Wrapf(
			err,
			"command %s failed with output: %s",
			color.YellowString(Line(cmd)),
			strings.Trim(out, "\n"),
		)
	}

	return
}

func (e *executor) run(ctx context.Context, out *bytes.Buffer, cmd *exec.Cmd) (string, error) {
	e.logger.WithField("args", cmd.Args).Debugf("executing command")

	if err := cmd.Start(); err != nil {
		return "", err
	}

	waitDone := make(chan struct{})
	ctxDone := ctx.Done()

	cancelSignal := os.Interrupt
	if s, ok := ctx.Value(CancelSignal).(os.Signal); ok {
		cancelSignal = s
	}

	go func() {
		select {
		case <-ctxDone:
			if cmd.Process != nil {
				e.logger.Infof("terminating running process...")
				cmd.Process.Signal(cancelSignal)
			}
		case <-waitDone:
		}
	}()

	err := cmd.Wait()

	close(waitDone)

	return out.String(), err
}

// Line returns the command line string for cmd.
func Line(cmd *exec.Cmd) string {
	return strings.Join(cmd.Args, " ")
}

// SetExecutor sets the default executor.
func SetExecutor(e Executor) {
	DefaultExecutor = e
}

// SetExecutorWithRestore sets the default executor and returns a function that
// restores the previously set executor. Can be used to temporarily mock out the
// executor in tests.
func SetExecutorWithRestore(e Executor) func() {
	prevExecutor := DefaultExecutor

	SetExecutor(e)

	return func() {
		DefaultExecutor = prevExecutor
	}
}
