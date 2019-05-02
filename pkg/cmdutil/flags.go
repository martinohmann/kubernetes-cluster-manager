package cmdutil

import (
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/cluster"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/provisioner"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/renderer"
	"github.com/spf13/cobra"
)

var (
	debug bool
	quiet bool
)

// AddGlobalFlags adds globally available flags to cmd.
func AddGlobalFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug output")
	cmd.PersistentFlags().BoolVar(&quiet, "quiet", false, "Disable log output. Ignored if --debug is set")
}

// AddConfigFlag adds the config flag to cmd.
func AddConfigFlag(cmd *cobra.Command) {
	cmd.Flags().String("config", "", "Config file path")
}

// BindProvisionerFlags binds flags to provisioner options.
func BindProvisionerFlags(cmd *cobra.Command, o *provisioner.Options) {
	cmd.Flags().IntVar(&o.Terraform.Parallelism, "terraform-parallelism", 0, "Number of parallel terraform resource operations")
}

// BindRendererFlags binds flags to renderer options.
func BindRendererFlags(cmd *cobra.Command, o *renderer.Options) {
	cmd.Flags().StringVar(&o.Helm.ChartsDir, "helm-charts-dir", "./charts", "Path to helm charts")
}

// BindManagerFlags binds flags to options.
func BindManagerFlags(cmd *cobra.Command, o *cluster.Options) {
	cmd.Flags().BoolVar(&o.DryRun, "dry-run", false, "Do not make any changes")
	cmd.Flags().StringVar(&o.ManifestsDir, "manifests-dir", "./manifests", "Path to rendered manifests")
	cmd.Flags().StringVar(&o.Deletions, "deletions", "deletions.yaml", `Deletions file path`)
	cmd.Flags().StringVar(&o.Values, "values", "values.yaml", `Values file path`)
}
