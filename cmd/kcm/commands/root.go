package commands

import (
	"os"
	"os/exec"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/command"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/config"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/infra"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/manifest"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/provisioner"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	cfg = &config.Config{}

	rootCmd = &cobra.Command{
		Use:          "kcm",
		Short:        "Kubernetes Cluster Manager",
		SilenceUsage: true,
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
	rootCmd.PersistentFlags().StringVarP(&cfg.Kubeconfig, "kubeconfig", "k", "", "Path to kubeconfig file")
	rootCmd.PersistentFlags().StringVarP(&cfg.Server, "server", "s", "", "Kubernetes API server address")
	rootCmd.PersistentFlags().StringVarP(&cfg.Token, "token", "t", "", "Bearer token for authentication to the Kubernetes API server")
	rootCmd.PersistentFlags().StringVarP(&cfg.Manifest, "manifest", "m", "", `Manifest file path (default: "manifest.yaml")`)
	rootCmd.PersistentFlags().StringVarP(&cfg.Deletions, "deletions", "d", "", `Deletions file path (default: "deletions.yaml")`)
	rootCmd.PersistentFlags().StringVarP(&cfg.WorkingDir, "working-dir", "w", workingDir, "Working directory")
	rootCmd.PersistentFlags().IntVar(&cfg.Terraform.Parallelism, "terraform-parallism", 1, "Number of parallel terraform resource operations")
	rootCmd.PersistentFlags().StringVar(&cfg.Helm.Values, "helm-values", "", `Values file path (default: "values.yaml")`)
	rootCmd.PersistentFlags().StringVar(&cfg.Helm.Chart, "helm-chart", "", `Path to cluster helm chart (default: "cluster")`)
}

func setupEnvironment() {
	cfg.ApplyDefaults()

	if err := os.Chdir(cfg.WorkingDir); err != nil {
		log.Fatal(err)
	}
}

func createProvisioner() *provisioner.Provisioner {
	if cfg.Debug {
		log.SetLevel(log.DebugLevel)
	}

	executor := command.NewExecutor()
	infraManager := infra.NewTerraformManager(cfg, executor)
	manifestRenderer := manifest.NewHelmRenderer(cfg, executor)

	return provisioner.NewClusterProvisioner(infraManager, manifestRenderer, executor)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			log.Println(err)
			os.Exit(exitErr.ExitCode())
		} else {
			log.Fatal(err)
		}
	}
}
