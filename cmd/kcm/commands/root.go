package commands

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/martinohmann/kubernetes-cluster-manager/infra"
	"github.com/martinohmann/kubernetes-cluster-manager/manifest"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/command"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/config"
	"github.com/martinohmann/kubernetes-cluster-manager/provisioner"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	cfg = &config.Config{}

	rootCmd = &cobra.Command{
		Use:           "kcm",
		Short:         "Kubernetes Cluster Manager",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
)

func init() {
	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	cobra.OnInitialize(setupEnvironment)

	rootCmd.PersistentFlags().BoolVar(&cfg.Debug, "debug", false, "Enable debug output")
	rootCmd.PersistentFlags().BoolVar(&cfg.DryRun, "dry-run", false, "Do not make any changes")
	rootCmd.PersistentFlags().BoolVar(&cfg.OnlyManifest, "only-manifest", false, "Only render manifest, skip infrastructure changes")
	rootCmd.PersistentFlags().StringVarP(&cfg.Cluster.Kubeconfig, "kubeconfig", "k", "", "Path to kubeconfig file")
	rootCmd.PersistentFlags().StringVar(&cfg.Cluster.Context, "kubeconfig-context", "", "Kubeconfig context")
	rootCmd.PersistentFlags().StringVarP(&cfg.Cluster.Server, "server", "s", "", "Kubernetes API server address")
	rootCmd.PersistentFlags().StringVarP(&cfg.Cluster.Token, "token", "t", "", "Bearer token for authentication to the Kubernetes API server")
	rootCmd.PersistentFlags().StringVarP(&cfg.Manifest, "manifest", "m", "", `Manifest file path (default: "manifest.yaml")`)
	rootCmd.PersistentFlags().StringVarP(&cfg.Deletions, "deletions", "d", "", `Deletions file path (default: "deletions.yaml")`)
	rootCmd.PersistentFlags().StringVar(&cfg.Values, "values", "", `Values file path (default: "values.yaml")`)
	rootCmd.PersistentFlags().StringVarP(&cfg.WorkingDir, "working-dir", "w", workingDir, "Working directory")
	rootCmd.PersistentFlags().IntVar(&cfg.Terraform.Parallelism, "terraform-parallism", 1, "Number of parallel terraform resource operations")
	rootCmd.PersistentFlags().StringVar(&cfg.Helm.Chart, "helm-chart", "", `Path to cluster helm chart (default: "cluster")`)
	rootCmd.PersistentFlags().StringVar(&cfg.InfraManager, "manager", "terraform", `Infrastructure manager to use`)
	rootCmd.PersistentFlags().StringVar(&cfg.ManifestRenderer, "renderer", "helm", `Manifest renderer to use`)
}

func setupEnvironment() {
	cfg.ApplyDefaults()

	if cfg.Debug {
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

	if err := os.Chdir(cfg.WorkingDir); err != nil {
		log.Fatal(err)
	}
}

func createProvisioner() (*provisioner.Provisioner, error) {
	executor := command.NewExecutor()
	infraManager, err := infra.CreateManager(cfg, executor)
	if err != nil {
		return nil, err
	}

	manifestRenderer, err := manifest.CreateRenderer(cfg, executor)
	if err != nil {
		return nil, err
	}

	p := provisioner.NewClusterProvisioner(infraManager, manifestRenderer, executor)

	return p, nil
}

// Execute executes the root command and prints eventual errors.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		code := 1
		cause := errors.Cause(err)

		if exitErr, ok := cause.(*exec.ExitError); ok {
			code = exitErr.ExitCode()
		}

		if cfg.Debug {
			log.Errorf("%+v", err)
		} else {
			log.Error(err)
		}

		os.Exit(code)
	}
}
