package cmdutil

import (
	"github.com/martinohmann/kubernetes-cluster-manager/infra"
	"github.com/martinohmann/kubernetes-cluster-manager/manifest"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kubernetes"
	"github.com/martinohmann/kubernetes-cluster-manager/provisioner"
	"github.com/spf13/cobra"
)

func AddGlobalDebugFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug output")
}

func AddConfigFlag(cmd *cobra.Command) {
	cmd.Flags().String("config", "", "Config file path")
}

func BindInfraManagerOptions(cmd *cobra.Command, o *infra.ManagerOptions) {
	cmd.Flags().IntVar(&o.Terraform.Parallelism, "terraform-parallelism", 1, "Number of parallel terraform resource operations")
}

func BindManifestRendererFlags(cmd *cobra.Command, o *manifest.RendererOptions) {
	cmd.Flags().StringVar(&o.Helm.Chart, "helm-chart", "./cluster", "Path to cluster helm chart")
}

func BindClusterFlags(cmd *cobra.Command, o *kubernetes.ClusterOptions) {
	cmd.Flags().StringVar(&o.Kubeconfig, "cluster-kubeconfig", "", "Path to kubeconfig file")
	cmd.Flags().StringVar(&o.Context, "cluster-context", "", "Kubeconfig context")
	cmd.Flags().StringVar(&o.Server, "cluster-server", "", "Kubernetes API server address")
	cmd.Flags().StringVar(&o.Token, "cluster-token", "", "Bearer token for authentication to the Kubernetes API server")
}

func BindProvisionerFlags(cmd *cobra.Command, o *provisioner.Options) {
	cmd.Flags().BoolVar(&o.DryRun, "dry-run", false, "Do not make any changes")
	cmd.Flags().BoolVar(&o.OnlyManifest, "only-manifest", false, "Only render manifest, skip infrastructure changes")
	cmd.Flags().StringVarP(&o.Manifest, "manifest", "m", "manifest.yaml", `Manifest file path`)
	cmd.Flags().StringVarP(&o.Deletions, "deletions", "d", "deletions.yaml", `Deletions file path`)
	cmd.Flags().StringVar(&o.Values, "values", "values.yaml", `Values file path`)
}
