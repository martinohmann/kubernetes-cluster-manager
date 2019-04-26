package cmdutil

import (
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kcm"
	"github.com/spf13/cobra"
)

var debug bool

func AddGlobalDebugFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug output")
}

func AddConfigFlag(cmd *cobra.Command) {
	cmd.Flags().String("config", "", "Config file path")
}

func BindProvisionerFlags(cmd *cobra.Command, o *kcm.ProvisionerOptions) {
	cmd.Flags().IntVar(&o.Terraform.Parallelism, "terraform-parallelism", 1, "Number of parallel terraform resource operations")
}

func BindRendererFlags(cmd *cobra.Command, o *kcm.RendererOptions) {
	cmd.Flags().StringVar(&o.Helm.Chart, "helm-chart", "./cluster", "Path to cluster helm chart")
}

func BindManagerFlags(cmd *cobra.Command, o *kcm.Options) {
	cmd.Flags().BoolVar(&o.DryRun, "dry-run", false, "Do not make any changes")
	cmd.Flags().BoolVar(&o.OnlyManifest, "only-manifest", false, "Only render manifest, skip infrastructure changes")
	cmd.Flags().StringVarP(&o.Manifest, "manifest", "m", "manifest.yaml", `Manifest file path`)
	cmd.Flags().StringVarP(&o.Deletions, "deletions", "d", "deletions.yaml", `Deletions file path`)
	cmd.Flags().StringVar(&o.Values, "values", "values.yaml", `Values file path`)
}
