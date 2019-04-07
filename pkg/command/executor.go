package command

import (
	"bufio"
	"bytes"
	"io"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Executor defines the interface for a command executor.
type Executor interface {
	// Run executes given command and returns its output.
	Run(*exec.Cmd) (string, error)
}

type executor func(*exec.Cmd) (string, error)

// NewExecutor creates a new command executor.
func NewExecutor() executor {
	return executor(Run)
}

// Run implements Executor.
func (e executor) Run(cmd *exec.Cmd) (string, error) {
	return Run(cmd)
}

// Run executes given command and returns its output. Will use the default
// *logrus.Logger to log cmd's stdout and stderr.
func Run(cmd *exec.Cmd) (string, error) {
	var out bytes.Buffer

	cmd.Stdout = io.MultiWriter(&out, logWriter(log.Info))
	cmd.Stderr = io.MultiWriter(&out, logWriter(log.Error))

	log.Debugf("Executing %s", strings.Join(cmd.Args, " "))

	err := cmd.Run()

	return out.String(), err
}

// logWriter wraps a logging function with an io.Writer
type logWriter func(args ...interface{})

// Write implements io.Writer.
func (w logWriter) Write(p []byte) (n int, err error) {
	s := bufio.NewScanner(bytes.NewReader(p))
	s.Split(bufio.ScanLines)

	for s.Scan() {
		w(s.Text())
	}

	return len(p), nil
}
