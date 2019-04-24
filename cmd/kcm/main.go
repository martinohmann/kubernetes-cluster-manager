package main

import (
	"os"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/cmd"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/cmdutil"
	log "github.com/sirupsen/logrus"
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
	logger := log.New()

	cmdutil.SetLogger(logger)
	cmdutil.AddGlobalDebugFlag(rootCmd)

	rootCmd.AddCommand(cmd.NewProvisionCommand(logger))
	rootCmd.AddCommand(cmd.NewDestroyCommand(logger))
	rootCmd.AddCommand(cmd.NewDumpConfigCommand(os.Stdout))
	rootCmd.AddCommand(cmd.NewVersionCommand(os.Stdout))

	cobra.OnInitialize(cmdutil.SetupLogger)
}

func main() {
	cmdutil.CheckErr(rootCmd.Execute())
}
