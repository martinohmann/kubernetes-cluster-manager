package main

import (
	"os"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/cmd"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/cmdutil"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:           "kcm",
		Short:         "Kubernetes Cluster Manager",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
)

func init() {
	cmdutil.AddGlobalFlags(rootCmd)

	rootCmd.AddCommand(cmd.NewProvisionCommand())
	rootCmd.AddCommand(cmd.NewDestroyCommand())
	rootCmd.AddCommand(cmd.NewManifestsCommand())
	rootCmd.AddCommand(cmd.NewDumpConfigCommand(os.Stdout))
	rootCmd.AddCommand(cmd.NewVersionCommand(os.Stdout))

	cobra.OnInitialize(cmdutil.ConfigureLogging)
}

func main() {
	cmdutil.CheckErr(rootCmd.Execute())
}
