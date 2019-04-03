package git

import (
	"bytes"
	"os/exec"

	"github.com/martinohmann/cluster-manager/pkg/executor"
)

func Diff(a, b string) (string, error) {
	var buf bytes.Buffer

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

	if err := executor.Execute(&buf, args...); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			// Exit code 1 means that there is a diff
			if exitErr.ExitCode() == 1 {
				return buf.String(), nil
			}
		}

		return "", err
	}

	return buf.String(), nil
}
