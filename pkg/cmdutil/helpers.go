package cmdutil

import (
	"os"
	"os/exec"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// CheckErr logs err and exits with a non-zero code. If err is of type
// *exec.ExitError, the exit code from the error will be used. If err is nil,
// CheckErr does nothing.
func CheckErr(err error) {
	if err == nil {
		return
	}

	code := 1
	cause := errors.Cause(err)

	if exitErr, ok := cause.(*exec.ExitError); ok {
		code = exitErr.ExitCode()
	}

	if debug {
		logger.Errorf("%+v", err)
	} else {
		logger.Error(err)
	}

	os.Exit(code)
}

// GetString retrieves a string flag from cmd.
func GetString(cmd *cobra.Command, flag string) string {
	v, err := cmd.Flags().GetString(flag)
	if err != nil {
		log.Fatal(err)
	}

	return v
}
