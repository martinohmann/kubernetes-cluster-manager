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
	cmd.Flags().IntVar(&o.Terraform.Parallelism, "terraform-parallelism", 0, "Number of parallel terraform resource operations")
}

func BindRendererFlags(cmd *cobra.Command, o *kcm.RendererOptions) {
	cmd.Flags().StringVar(&o.Helm.ChartsDir, "helm-charts-dir", "./charts", "Path to helm charts")
}

func BindManagerFlags(cmd *cobra.Command, o *kcm.Options) {
	cmd.Flags().BoolVar(&o.DryRun, "dry-run", false, "Do not make any changes")
	cmd.Flags().StringVar(&o.ManifestsDir, "manifests-dir", "./manifests", "Path to rendered manifests")
	cmd.Flags().StringVar(&o.Deletions, "deletions", "deletions.yaml", `Deletions file path`)
	cmd.Flags().StringVar(&o.Values, "values", "values.yaml", `Values file path`)
}
