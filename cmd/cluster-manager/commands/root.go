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
		Use: "cluster-manager",
	}
)

func init() {
	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	cobra.OnInitialize(cfg.ApplyDefaults)

	rootCmd.PersistentFlags().BoolVar(&cfg.DryRun, "dry-run", false, "Do not make any changes")
	rootCmd.PersistentFlags().StringVar(&cfg.Kubeconfig, "kubeconfig", "", "Kubeconfig")
	rootCmd.PersistentFlags().StringVar(&cfg.Manifest, "manifest", "", "Manifest file path (Default: manifest.yaml)")
	rootCmd.PersistentFlags().StringVar(&cfg.Deletions, "deletions", "", "Deletions file path (Default: deletions.yaml)")
	rootCmd.PersistentFlags().StringVar(&cfg.WorkingDir, "working-dir", workingDir, "Working directory")
	rootCmd.PersistentFlags().BoolVar(&cfg.Terraform.AutoApprove, "terraform-auto-approve", false, "Automatically approve terraform changes")
	rootCmd.PersistentFlags().IntVar(&cfg.Terraform.Parallelism, "terraform-parallism", 1, "Number of parallel terraform resource operations")
	rootCmd.PersistentFlags().StringVar(&cfg.Helm.Values, "helm-values", "", "Values file path (Default: values.yaml)")
	rootCmd.PersistentFlags().StringVar(&cfg.Helm.Chart, "helm-chart", "", "Path to cluster helm chart (Default: cluster)")
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
