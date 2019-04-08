package commands

import (
	"github.com/spf13/cobra"
)

var (
	provisionCmd = &cobra.Command{
		Use:   "provision",
		Short: "Provisions a cluster",
		Long: "Provisions a cluster by creating or updating its infrastructure\n" +
			"resources and afterwards applying Kubernetes manifests.",
		RunE: provision,
	}
)

func init() {
	rootCmd.AddCommand(provisionCmd)
}

func provision(cmd *cobra.Command, args []string) error {
	provisioner := createProvisioner()

	return provisioner.Provision(cfg)
}
