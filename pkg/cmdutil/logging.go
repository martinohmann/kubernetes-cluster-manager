package cmdutil

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"
)

var logger = log.StandardLogger()

// ConfigureLogger configures l based on the values of parsed cli flags.
func ConfigureLogger(l *log.Logger) {
	logger = l

	if quiet && !debug {
		logger.Out = ioutil.Discard
	}

	if !debug {
		return
	}

	logger.SetLevel(log.DebugLevel)
	logger.SetReportCaller(true)
	logger.SetFormatter(&log.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			pkg := "github.com/martinohmann/kubernetes-cluster-manager/"
			repopath := fmt.Sprintf("%s/src/%s", os.Getenv("GOPATH"), pkg)
			filename := strings.Replace(f.File, repopath, "", -1)
			function := strings.Replace(f.Function, pkg, "", -1)
			return fmt.Sprintf("%s()", function), fmt.Sprintf("%s:%d", filename, f.Line)
		},
	})
}
