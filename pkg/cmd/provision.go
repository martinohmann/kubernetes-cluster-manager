package cmd

import (
	"os"

	"github.com/fatih/color"
	"github.com/imdario/mergo"
	"github.com/martinohmann/kubernetes-cluster-manager/infra"
	"github.com/martinohmann/kubernetes-cluster-manager/manifest"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/cmdutil"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/command"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/file"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kubernetes"
	"github.com/martinohmann/kubernetes-cluster-manager/provisioner"
	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type Options struct {
	Manager    string `json:"manager,omitempty" yaml:"manager,omitempty"`
	Renderer   string `json:"renderer,omitempty" yaml:"renderer,omitempty"`
	WorkingDir string `json:"workingDir,omitempty" yaml:"workingDir,omitempty"`

	ProvisionerOptions      provisioner.Options       `json:"provisioner,omitempty" yaml:"provisioner,omitempty"`
	ClusterOptions          kubernetes.ClusterOptions `json:"cluster,omitempty" yaml:"cluster,omitempty"`
	InfraManagerOptions     infra.ManagerOptions      `json:"infraManager,omitempty" yaml:"infraManager,omitempty"`
	ManifestRendererOptions manifest.RendererOptions  `json:"manifestRenderer,omitempty" yaml:"manifestRenderer,omitempty"`

	destroy bool
	logger  *log.Logger
}

func NewProvisionCommand(l *log.Logger) *cobra.Command {
	o := &Options{
		destroy: false,
		logger:  l,
	}

	cmd := &cobra.Command{
		Use:   "provision",
		Short: "Provisions a cluster",
		Long: "Provisions a cluster by creating or updating its infrastructure\n" +
			"resources and afterwards applying Kubernetes manifests.",
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(cmd))
			cmdutil.CheckErr(o.Run())
		},
	}

	o.AddFlags(cmd)

	return cmd
}

func (o *Options) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.Manager, "manager", "terraform", `Infrastructure manager to use`)
	cmd.Flags().StringVar(&o.Renderer, "renderer", "helm", `Manifest renderer to use`)
	cmd.Flags().StringVarP(&o.WorkingDir, "working-dir", "w", "", "Working directory")

	cmdutil.AddConfigFlag(cmd)
	cmdutil.BindProvisionerFlags(cmd, &o.ProvisionerOptions)
	cmdutil.BindClusterFlags(cmd, &o.ClusterOptions)
	cmdutil.BindInfraManagerOptions(cmd, &o.InfraManagerOptions)
	cmdutil.BindManifestRendererFlags(cmd, &o.ManifestRendererOptions)
}

func (o *Options) Complete(cmd *cobra.Command) error {
	var err error

	if config := cmdutil.GetString(cmd, "config"); config != "" {
		if err = o.MergeConfig(config); err != nil {
			return err
		}

		o.logger.Infof("Using config %s, config values take precedence over command line flags", color.YellowString(config))

	}

	o.WorkingDir, err = homedir.Expand(o.WorkingDir)
	if err != nil {
		return err
	}

	return nil
}

func (o *Options) Run() error {
	if o.WorkingDir != "" {
		if err := os.Chdir(o.WorkingDir); err != nil {
			return err
		}
	}

	p, err := o.createProvisioner()
	if err != nil {
		return err
	}

	if o.destroy {
		return p.Destroy(&o.ProvisionerOptions)
	}

	return p.Provision(&o.ProvisionerOptions)
}

func (o *Options) MergeConfig(filename string) error {
	opts := &Options{}

	if err := file.LoadYAML(filename, opts); err != nil {
		return err
	}

	return mergo.Merge(o, opts, mergo.WithOverride)
}

func (o *Options) createProvisioner() (*provisioner.Provisioner, error) {
	executor := command.NewExecutor(o.logger)
	infraManager, err := infra.CreateManager(o.Manager, &o.InfraManagerOptions, executor)
	if err != nil {
		return nil, err
	}

	manifestRenderer, err := manifest.CreateRenderer(o.Renderer, &o.ManifestRendererOptions, executor)
	if err != nil {
		return nil, err
	}

	p := provisioner.NewClusterProvisioner(
		&o.ClusterOptions,
		infraManager,
		manifestRenderer,
		executor,
		o.logger,
	)

	return p, nil
}
