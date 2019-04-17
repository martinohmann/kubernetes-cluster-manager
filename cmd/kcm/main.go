package main

import (
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
	cobra.OnInitialize(cmdutil.SetupLogger)

	cmdutil.AddGlobalDebugFlag(rootCmd)

	rootCmd.AddCommand(cmd.NewProvisionCommand())
	rootCmd.AddCommand(cmd.NewDestroyCommand())
	rootCmd.AddCommand(cmd.NewDumpConfigCommand())
	rootCmd.AddCommand(cmd.NewVersionCommand())
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
