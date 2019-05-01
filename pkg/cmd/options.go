package cmd

import (
	"os"

	"github.com/fatih/color"
	"github.com/imdario/mergo"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/cluster"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/cmdutil"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/command"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/credentials"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/file"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kcm"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/provisioner"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/renderer"
	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type ClusterOptions struct {
	Server     string `json:"server" yaml:"server"`
	Token      string `json:"token" yaml:"token"`
	Kubeconfig string `json:"kubeconfig" yaml:"kubeconfig"`
	Context    string `json:"context" yaml:"context"`
}

type Options struct {
	Provisioner string `json:"provisioner,omitempty" yaml:"provisioner,omitempty"`
	Renderer    string `json:"renderer,omitempty" yaml:"renderer,omitempty"`
	WorkingDir  string `json:"workingDir,omitempty" yaml:"workingDir,omitempty"`

	ClusterOptions     ClusterOptions         `json:"clusterCredentials,omitempty" yaml:"clusterOptions,omitempty"`
	ManagerOptions     kcm.Options            `json:"managerOptions,omitempty" yaml:"managerOptions,omitempty"`
	ProvisionerOptions kcm.ProvisionerOptions `json:"provisionerOptions,omitempty" yaml:"provisionerOptions,omitempty"`
	RendererOptions    kcm.RendererOptions    `json:"rendererOptions,omitempty" yaml:"rendererOptions,omitempty"`

	logger *log.Logger
}

func (o *Options) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.Provisioner, "provisioner", "", `Infrastructure provisioner to use`)
	cmd.Flags().StringVar(&o.Renderer, "renderer", "helm", `Manifest renderer to use`)
	cmd.Flags().StringVarP(&o.WorkingDir, "working-dir", "w", "", "Working directory")

	cmd.Flags().StringVar(&o.ClusterOptions.Kubeconfig, "cluster-kubeconfig", "", "Path to kubeconfig file")
	cmd.Flags().StringVar(&o.ClusterOptions.Context, "cluster-context", "", "Kubeconfig context")
	cmd.Flags().StringVar(&o.ClusterOptions.Server, "cluster-server", "", "Kubernetes API server address")
	cmd.Flags().StringVar(&o.ClusterOptions.Token, "cluster-token", "", "Bearer token for authentication to the Kubernetes API server")

	cmdutil.AddConfigFlag(cmd)
	cmdutil.BindManagerFlags(cmd, &o.ManagerOptions)
	cmdutil.BindProvisionerFlags(cmd, &o.ProvisionerOptions)
	cmdutil.BindRendererFlags(cmd, &o.RendererOptions)
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
	if o.Provisioner == "" {
		o.Provisioner = "null"
	}

	executor := command.NewExecutor(o.logger)

	command.SetExecutor(executor)

	return err
}

func (o *Options) Run(exec func(kcm.ClusterManager, *kcm.Options) error) error {
	if o.WorkingDir != "" {
		o.logger.Infof("Switching working dir to %s", o.WorkingDir)
		if err := os.Chdir(o.WorkingDir); err != nil {
			return err
		}
	}

	m, err := o.createManager()
	if err != nil {
		return err
	}

	return exec(m, &o.ManagerOptions)
}

func (o *Options) MergeConfig(filename string) error {
	opts := &Options{}

	if err := file.ReadYAML(filename, opts); err != nil {
		return err
	}

	return mergo.Merge(o, opts, mergo.WithOverride)
}

func (o *Options) createManager() (kcm.ClusterManager, error) {
	infraProvisioner, err := provisioner.Create(o.Provisioner, &o.ProvisionerOptions)
	if err != nil {
		return nil, err
	}

	manifestRenderer, err := renderer.Create(o.Renderer, &o.RendererOptions)
	if err != nil {
		return nil, err
	}

	var credentialSource kcm.CredentialSource
	if o.ClusterOptions.Kubeconfig == "" && o.ClusterOptions.Context == "" &&
		(o.ClusterOptions.Server == "" || o.ClusterOptions.Token == "") {
		credentialSource = credentials.NewProvisionerSource(infraProvisioner)
	} else {
		credentialSource = credentials.NewStaticCredentials(&kcm.Credentials{
			Server:     o.ClusterOptions.Server,
			Token:      o.ClusterOptions.Token,
			Kubeconfig: o.ClusterOptions.Kubeconfig,
			Context:    o.ClusterOptions.Context,
		})
	}

	return cluster.NewManager(
		credentialSource,
		infraProvisioner,
		manifestRenderer,
		o.logger,
	), nil
}
