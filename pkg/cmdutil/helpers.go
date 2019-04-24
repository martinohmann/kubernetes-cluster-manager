package cmdutil

import (
	"os"
	"os/exec"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

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

func GetBool(cmd *cobra.Command, flag string) bool {
	v, err := cmd.Flags().GetBool(flag)
	if err != nil {
		log.Fatal(err)
	}

	return v
}

func GetString(cmd *cobra.Command, flag string) string {
	v, err := cmd.Flags().GetString(flag)
	if err != nil {
		log.Fatal(err)
	}

	return v
}

func GetInt(cmd *cobra.Command, flag string) int {
	v, err := cmd.Flags().GetInt(flag)
	if err != nil {
		log.Fatal(err)
	}

	return v
}
