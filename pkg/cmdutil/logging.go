package cmdutil

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"
)

// ConfigureLogging configures the standard logger based on the values of parsed
// cli flags.
func ConfigureLogging() {
	if quiet && !debug {
		log.SetOutput(ioutil.Discard)
	}

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
