package git

import (
	"os/exec"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/command"
)

// DiffTool is a utility for working with file changes.
type DiffTool struct {
	DiffOnly bool
}

// Apply makes a diff of the FileChanges and returns it. Will apply the changes
// to the source file unless DiffTool is configured with DiffOnly.
func (t *DiffTool) Apply(changes *FileChanges) (string, error) {
	diff, err := Diff(changes.filename, changes.tmpf.Name())
	if err != nil {
		return diff, err
	}

	if t.DiffOnly {
		return diff, nil
	}

	return diff, changes.Apply()
}

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
