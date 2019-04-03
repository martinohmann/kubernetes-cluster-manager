package commands

import (
	"os"

	"github.com/martinohmann/cluster-manager/pkg/infra/terraform"
	"github.com/martinohmann/cluster-manager/pkg/manifest/helm"
	"github.com/martinohmann/cluster-manager/pkg/provisioner"

	"github.com/spf13/cobra"
)

var (
	provisionCmd = &cobra.Command{
		Use:  "provision",
		RunE: provision,
	}
)

func init() {
	rootCmd.AddCommand(provisionCmd)
}

func provision(cmd *cobra.Command, args []string) (err error) {
	if err = os.Chdir(cfg.WorkingDir); err != nil {
		return
	}

	writer := os.Stdout

	m := terraform.NewInfraManager(cfg, writer)
	r := helm.NewManifestRenderer(cfg)
	p := provisioner.NewClusterProvisioner(m, r, writer)

	if err = p.Provision(cfg, nil); err != nil {
		return
	}

	return
}
