package cmd

import (
	"context"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/cluster"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewDestroyCommand() *cobra.Command {
	o := &Options{}

	cmd := &cobra.Command{
		Use:   "destroy",
		Short: "Destroys a cluster",
		Long: "Destroys a Kubernetes cluster by first deleting all resources defined\n" +
			"in the manifest and afterwards deleting all infrastructure resources.",
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(cmd))
			cmdutil.CheckErr(o.Run(func(ctx context.Context, m *cluster.Manager, o *cluster.Options) error {
				return m.Destroy(ctx, o)
			}))
		},
	}

	o.AddFlags(cmd)

	cmd.Flags().BoolVar(&o.ManagerOptions.SkipManifests, "skip-manifests", false, "Skip processing kubernetes manifests")

	return cmd
}
