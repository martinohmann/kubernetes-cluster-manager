package cmd

import (
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/cmdutil"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kcm"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewManifestsCommand(l *log.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "manifests",
		Aliases: []string{"manifest"},
		Short:   "Perform manifest actions",
	}

	cmd.AddCommand(newApplyCommand(l))
	cmd.AddCommand(newDeleteCommand(l))

	return cmd
}

func newApplyCommand(l *log.Logger) *cobra.Command {
	o := &Options{logger: l}

	cmd := &cobra.Command{
		Use:   "apply",
		Short: "Applies manifests to a cluster",
		Long:  "Renders manifests and applies them to a cluster.",
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(cmd))
			cmdutil.CheckErr(o.Run(func(m kcm.ClusterManager, o *kcm.Options) error {
				return m.ApplyManifests(o)
			}))
		},
	}

	o.AddFlags(cmd)

	return cmd
}

func newDeleteCommand(l *log.Logger) *cobra.Command {
	o := &Options{logger: l}

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Deletes manifests from a cluster",
		Long:  "Renders manifests and deletes them from a cluster.",
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(cmd))
			cmdutil.CheckErr(o.Run(func(m kcm.ClusterManager, o *kcm.Options) error {
				return m.DeleteManifests(o)
			}))
		},
	}

	o.AddFlags(cmd)

	return cmd
}
