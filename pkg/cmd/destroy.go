package cmd

import (
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/cmdutil"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewDestroyCommand(l *log.Logger) *cobra.Command {
	o := &Options{
		destroy: true,
		logger:  l,
	}

	cmd := &cobra.Command{
		Use:   "destroy",
		Short: "Destroys a cluster",
		Long: "Destroys a Kubernetes cluster by first deleting all resources defined\n" +
			"in the manifest and afterwards deleting all infrastructure resources.",
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(cmd))
			cmdutil.CheckErr(o.Run())
		},
	}

	o.AddFlags(cmd)

	return cmd
}
