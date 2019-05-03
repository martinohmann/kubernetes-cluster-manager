package cmd

import (
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/cluster"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/cmdutil"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewProvisionCommand(l *log.Logger) *cobra.Command {
	o := &Options{logger: l}

	cmd := &cobra.Command{
		Use:   "provision",
		Short: "Provisions a cluster",
		Long: "Provisions a cluster by creating or updating its infrastructure\n" +
			"resources and afterwards applying Kubernetes manifests.",
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(cmd))
			cmdutil.CheckErr(o.Run(func(m *cluster.Manager, o *cluster.Options) error {
				return m.Provision(o)
			}))
		},
	}

	o.AddFlags(cmd)

	cmd.Flags().BoolVar(&o.ManagerOptions.SkipManifests, "skip-manifests", false, "Skip processing kubernetes manifests")
	cmd.Flags().BoolVar(&o.ManagerOptions.AllManifests, "all-manifests", false, "Apply all manifests, even unchanged")

	return cmd
}
