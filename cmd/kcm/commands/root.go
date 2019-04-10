package commands

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/command"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/config"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/factories"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/provisioner"
	"github.com/sirupsen/logrus"
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

	if err := os.Chdir(cfg.WorkingDir); err != nil {
		log.Fatal(err)
	}

	if cfg.Debug {
		log.SetLevel(log.DebugLevel)
		log.SetReportCaller(true)
		log.SetFormatter(&logrus.TextFormatter{
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				pkg := "github.com/martinohmann/kubernetes-cluster-manager/"
				repopath := fmt.Sprintf("%s/src/%s", os.Getenv("GOPATH"), pkg)
				filename := strings.Replace(f.File, repopath, "", -1)
				function := strings.Replace(f.Function, pkg, "", -1)
				return fmt.Sprintf("%s()", function), fmt.Sprintf("%s:%d", filename, f.Line)
			},
		})
	}
}

func createProvisioner() (*provisioner.Provisioner, error) {
	executor := command.NewExecutor()
	infraManager, err := factories.CreateInfraManager(cfg, executor)
	if err != nil {
		return nil, err
	}

	manifestRenderer, err := factories.CreateManifestRenderer(cfg, executor)
	if err != nil {
		return nil, err
	}

	p := provisioner.NewClusterProvisioner(infraManager, manifestRenderer, executor)

	return p, nil
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			log.Printf("%+v", err)
			os.Exit(exitErr.ExitCode())
		} else {
			log.Fatalf("%+v", err)
		}
	}
}
