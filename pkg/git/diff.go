package git

import (
	"os/exec"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/command"
)

// Diff runs `git diff` on files a and b. Will return the diff and an error if
// `git diff` fails.
func Diff(a, b string) (out string, err error) {
	args := []string{
		"git",
		"--no-pager",
		"diff",
		"--color=always",
		"--no-index",
		"--exit-code",
		a,
		b,
	}

	cmd := exec.Command(args[0], args[1:]...)

	if out, err = command.RunSilently(cmd); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			// Exit code 1 means that there is a diff
			if exitErr.ExitCode() == 1 {
				err = nil
			}
		}
	}

	return out, err
}

// DiffFileChanges creates a diff for the file changes and returns it.
func DiffFileChanges(changes *FileChanges) (string, error) {
	return Diff(changes.filename, changes.tmpf.Name())
}
