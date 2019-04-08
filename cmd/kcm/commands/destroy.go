package commands

import (
	"github.com/spf13/cobra"
)

var (
	destroyCmd = &cobra.Command{
		Use:   "destroy",
		Short: "Destroys a cluster",
		Long: "Destroys a Kubernetes cluster by first deleting all resources defined\n" +
			"in the manifest and afterwards deleting all infrastructure resources.",
		RunE: destroy,
	}
)

func init() {
	rootCmd.AddCommand(destroyCmd)
}

func destroy(cmd *cobra.Command, args []string) error {
	provisioner := createProvisioner()

	return provisioner.Destroy(cfg)
}
