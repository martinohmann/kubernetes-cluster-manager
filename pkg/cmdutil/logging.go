package cmdutil

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var debug = false

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
		log.Errorf("%+v", err)
	} else {
		log.Error(err)
	}

	os.Exit(code)
}

func SetupLogger() {
	if !debug {
		return
	}

	log.SetLevel(log.DebugLevel)
	log.SetReportCaller(true)
	log.SetFormatter(&log.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			pkg := "github.com/martinohmann/kubernetes-cluster-manager/"
			repopath := fmt.Sprintf("%s/src/%s", os.Getenv("GOPATH"), pkg)
			filename := strings.Replace(f.File, repopath, "", -1)
			function := strings.Replace(f.Function, pkg, "", -1)
			return fmt.Sprintf("%s()", function), fmt.Sprintf("%s:%d", filename, f.Line)
		},
	})
}
