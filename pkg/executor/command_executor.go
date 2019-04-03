package executor

import (
	"io"
	"os/exec"
)

func Execute(out io.Writer, cmdArgs ...string) error {
	return Pipe(nil, out, cmdArgs...)
}

func Pipe(in io.Reader, out io.Writer, cmdArgs ...string) error {
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	cmd.Stdin = in
	cmd.Stdout = out
	cmd.Stderr = out

	return cmd.Run()
}
