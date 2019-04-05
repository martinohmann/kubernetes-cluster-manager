package commands

import (
	"os"

	"github.com/martinohmann/cluster-manager/pkg/infra"
	"github.com/martinohmann/cluster-manager/pkg/manifest"
	"github.com/martinohmann/cluster-manager/pkg/provisioner"

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
	writer := os.Stdout

	m := infra.NewTerraformManager(cfg, writer)
	r := manifest.NewHelmRenderer(cfg)
	p := provisioner.NewClusterProvisioner(m, r, writer)

	return p.Provision(cfg)
}
