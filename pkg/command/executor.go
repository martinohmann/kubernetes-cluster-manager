package command

import (
	"bufio"
	"bytes"
	"io"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"
)

// Executor defines the interface for a command executor.
type Executor interface {
	// Run executes given command and returns its output.
	Run(*exec.Cmd) (string, error)

	// RunSilently executes the given command and returns its output. Will not
	// write command output to stdout or stderr.
	RunSilently(*exec.Cmd) (string, error)
}

// DefaultExecutor is the default executor used in the package level Run and
// RunSilently funcs.
var DefaultExecutor = NewExecutor()

// Run runs a command using the default executor.
func Run(cmd *exec.Cmd) (string, error) {
	return DefaultExecutor.Run(cmd)
}

// RunSilently runs cmd silently using the default executor.
func RunSilently(cmd *exec.Cmd) (string, error) {
	return DefaultExecutor.RunSilently(cmd)
}

type executor struct {
	logger *log.Logger
}

// NewExecutor creates a new command executor. Optionally accepts a logger for
// logging command output. If not provided logrus.StandardLogger() will be used.
func NewExecutor(logger ...*log.Logger) Executor {
	l := log.StandardLogger()
	if len(logger) > 0 && logger[0] != nil {
		l = logger[0]
	}

	return &executor{logger: l}
}

// Run implements Run from Executor interface.
func (e *executor) Run(cmd *exec.Cmd) (string, error) {
	var out bytes.Buffer

	cmd.Stdout = io.MultiWriter(&out, newLogWriter(cmd, e.logger.Info))
	cmd.Stderr = io.MultiWriter(&out, newLogWriter(cmd, e.logger.Error))

	return e.run(&out, cmd)
}

// RunSilently implements RunSilently from Executor interface.
func (e *executor) RunSilently(cmd *exec.Cmd) (string, error) {
	var out bytes.Buffer

	cmd.Stdout = &out
	cmd.Stderr = &out

	return e.run(&out, cmd)
}

func (e *executor) run(out *bytes.Buffer, cmd *exec.Cmd) (string, error) {
	e.logger.Debugf("Executing %s", color.YellowString(commandLine(cmd)))

	err := cmd.Run()

	return out.String(), err
}

func newLogWriter(cmd *exec.Cmd, f func(...interface{})) logWriter {
	return logWriter{
		prefix: color.BlueString("[%s] ", cmd.Args[0]),
		f:      f,
	}
}

// logWriter wraps a logging function with an io.Writer
type logWriter struct {
	prefix string
	f      func(args ...interface{})
}

// Write implements io.Writer.
func (w logWriter) Write(p []byte) (n int, err error) {
	s := bufio.NewScanner(bytes.NewReader(p))
	s.Split(bufio.ScanLines)

	for s.Scan() {
		w.f(w.prefix + s.Text())
	}

	return len(p), nil
}

func commandLine(cmd *exec.Cmd) string {
	return strings.Join(cmd.Args, " ")
}
