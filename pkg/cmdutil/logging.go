package cmdutil

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/log"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"
)

// ConfigureLogging configures the standard logger based on the values of parsed
// cli flags.
func ConfigureLogging() {
	formatter := &log.ContextPrefixFormatter{&logrus.TextFormatter{}}

	logrus.SetFormatter(formatter)

	if quiet && !debug {
		logrus.SetOutput(ioutil.Discard)
	} else {
		logrus.SetOutput(os.Stdout)

		if !terminal.IsTerminal(int(os.Stdout.Fd())) {
			color.NoColor = true
		}
	}

	if debug {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.SetReportCaller(true)

		formatter.FullTimestamp = true
		formatter.CallerPrettyfier = func(f *runtime.Frame) (string, string) {
			pkg := "github.com/martinohmann/kubernetes-cluster-manager/"
			repopath := fmt.Sprintf("%s/src/%s", os.Getenv("GOPATH"), pkg)
			filename := strings.Replace(f.File, repopath, "", -1)
			function := strings.Replace(f.Function, pkg, "", -1)
			return fmt.Sprintf("%s()", function), fmt.Sprintf("%s:%d", filename, f.Line)
		}
	}
}
