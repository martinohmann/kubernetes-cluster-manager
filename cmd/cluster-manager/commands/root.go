package commands

import (
	"log"
	"os"
	"os/exec"

	"github.com/martinohmann/cluster-manager/pkg/config"
	"github.com/spf13/cobra"
)

var (
	cfg = &config.Config{}

	rootCmd = &cobra.Command{
		Use:          "cluster-manager",
		SilenceUsage: true,
	}
)

func init() {
	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	cobra.OnInitialize(setupEnvironment)

	rootCmd.PersistentFlags().BoolVar(&cfg.DryRun, "dry-run", false, "Do not make any changes")
	rootCmd.PersistentFlags().BoolVar(&cfg.OnlyManifest, "only-manifest", false, "Only render manifest, skip infrastructure changes")
	rootCmd.PersistentFlags().StringVarP(&cfg.Kubeconfig, "kubeconfig", "k", "", "Path to kubeconfig file")
	rootCmd.PersistentFlags().StringVarP(&cfg.Manifest, "manifest", "m", "", `Manifest file path (default: "manifest.yaml")`)
	rootCmd.PersistentFlags().StringVarP(&cfg.Deletions, "deletions", "d", "", `Deletions file path (default: "deletions.yaml")`)
	rootCmd.PersistentFlags().StringVarP(&cfg.WorkingDir, "working-dir", "w", workingDir, "Working directory")
	rootCmd.PersistentFlags().BoolVar(&cfg.Terraform.AutoApprove, "terraform-auto-approve", false, "Automatically approve terraform changes")
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
