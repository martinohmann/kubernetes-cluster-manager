package cmd

import (
	"context"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/cluster"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewManifestsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "manifests",
		Aliases: []string{"manifest"},
		Short:   "Perform manifest actions",
	}

	cmd.AddCommand(newApplyCommand())
	cmd.AddCommand(newDeleteCommand())

	return cmd
}

func newApplyCommand() *cobra.Command {
	o := &Options{}

	cmd := &cobra.Command{
		Use:   "apply",
		Short: "Applies manifests to a cluster",
		Long:  "Renders manifests and applies them to a cluster.",
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(cmd))
			cmdutil.CheckErr(o.Run(func(ctx context.Context, m *cluster.Manager, o *cluster.Options) error {
				return m.ApplyManifests(ctx, o)
			}))
		},
	}

	o.AddFlags(cmd)

	cmd.Flags().BoolVar(&o.ManagerOptions.AllManifests, "all-manifests", false, "Apply all manifests, even unchanged")
	cmd.Flags().BoolVar(&o.ManagerOptions.NoHooks, "no-hooks", false, "Skip executing hooks")

	return cmd
}

func newDeleteCommand() *cobra.Command {
	o := &Options{}

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Deletes manifests from a cluster",
		Long:  "Renders manifests and deletes them from a cluster.",
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(cmd))
			cmdutil.CheckErr(o.Run(func(ctx context.Context, m *cluster.Manager, o *cluster.Options) error {
				return m.DeleteManifests(ctx, o)
			}))
		},
	}

	o.AddFlags(cmd)

	cmd.Flags().BoolVar(&o.ManagerOptions.AllManifests, "all-manifests", false, "Attempt to delete all manifests, even the ones already absent")

	return cmd
}
