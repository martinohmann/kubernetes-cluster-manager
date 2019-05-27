package cmdutil

import (
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/cluster"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/provisioner"
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
	cmd.Flags().IntVar(&o.Parallelism, "parallelism", 0, "Number of parallel provisioner resource operations")
}

// BindManagerFlags binds flags to options.
func BindManagerFlags(cmd *cobra.Command, o *cluster.Options) {
	cmd.Flags().BoolVar(&o.DryRun, "dry-run", false, "Do not make any changes")
	cmd.Flags().StringVar(&o.ManifestsDir, "manifests-dir", "./manifests", "Path to rendered manifests")
	cmd.Flags().StringVar(&o.TemplatesDir, "templates-dir", "./templates", "Path to components containing manifest templates")
	cmd.Flags().StringVar(&o.Values, "values", "values.yaml", `Values file path`)
	cmd.Flags().BoolVar(&o.NoSave, "no-save", false, "Do not save file changes")
}
