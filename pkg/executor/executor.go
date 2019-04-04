package executor

import (
	"io"
	"os/exec"
)

var DefaultExecutor = &ShellExecutor{}

type Interface interface {
	Execute(io.Writer, ...string) error
	Pipe(io.Reader, io.Writer, ...string) error
}

type ShellExecutor struct{}

func (e *ShellExecutor) Execute(out io.Writer, cmdArgs ...string) error {
	return e.Pipe(nil, out, cmdArgs...)
}

func (e *ShellExecutor) Pipe(in io.Reader, out io.Writer, cmdArgs ...string) error {
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	cmd.Stdin = in
	cmd.Stdout = out
	cmd.Stderr = out

	return cmd.Run()
}

func Execute(out io.Writer, cmdArgs ...string) error {
	return DefaultExecutor.Execute(out, cmdArgs...)
}

func Pipe(in io.Reader, out io.Writer, cmdArgs ...string) error {
	return DefaultExecutor.Pipe(in, out, cmdArgs...)
}
