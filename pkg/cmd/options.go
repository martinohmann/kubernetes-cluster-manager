package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/fatih/color"
	"github.com/imdario/mergo"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/cluster"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/cmdutil"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/credentials"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/file"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/provisioner"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/renderer"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type Options struct {
	Provisioner string `json:"provisioner,omitempty" yaml:"provisioner,omitempty"`
	Renderer    string `json:"renderer,omitempty" yaml:"renderer,omitempty"`
	WorkingDir  string `json:"workingDir,omitempty" yaml:"workingDir,omitempty"`

	Credentials        credentials.Credentials `json:"credentials,omitempty" yaml:"credentials,omitempty"`
	ManagerOptions     cluster.Options         `json:"managerOptions,omitempty" yaml:"managerOptions,omitempty"`
	ProvisionerOptions provisioner.Options     `json:"provisionerOptions,omitempty" yaml:"provisionerOptions,omitempty"`
	RendererOptions    renderer.Options        `json:"rendererOptions,omitempty" yaml:"rendererOptions,omitempty"`
}

func (o *Options) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.Provisioner, "provisioner", "", `Infrastructure provisioner to use`)
	cmd.Flags().StringVar(&o.Renderer, "renderer", "helm", `Manifest renderer to use`)
	cmd.Flags().StringVarP(&o.WorkingDir, "working-dir", "w", "", "Working directory")

	cmd.Flags().StringVar(&o.Credentials.Kubeconfig, "cluster-kubeconfig", "", "Path to kubeconfig file")
	cmd.Flags().StringVar(&o.Credentials.Context, "cluster-context", "", "Kubeconfig context")
	cmd.Flags().StringVar(&o.Credentials.Server, "cluster-server", "", "Kubernetes API server address")
	cmd.Flags().StringVar(&o.Credentials.Token, "cluster-token", "", "Bearer token for authentication to the Kubernetes API server")

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

		log.Infof("Using config %s, config values take precedence over command line flags", color.YellowString(config))
	}

	o.WorkingDir, err = homedir.Expand(o.WorkingDir)
	if o.Provisioner == "" {
		o.Provisioner = "null"
	}

	if o.Renderer == "" {
		o.Renderer = "null"
	}

	return err
}

func (o *Options) Run(exec func(context.Context, *cluster.Manager, *cluster.Options) error) error {
	if o.WorkingDir != "" {
		log.Infof("Switching working dir to %s", o.WorkingDir)
		if err := os.Chdir(o.WorkingDir); err != nil {
			return err
		}
	}

	m, err := o.createManager()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan struct{}, 1)
	signalChan := make(chan os.Signal, 2)

	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(signalChan)

	go func() {
		for {
			select {
			case s := <-signalChan:
				log.Infof("Received signal %s, cleaning up...", s)
				cancel()
			case <-done:
				return
			}
		}
	}()

	err = exec(ctx, m, &o.ManagerOptions)

	close(done)

	return err
}

func (o *Options) MergeConfig(filename string) error {
	opts := &Options{}

	if err := file.ReadYAML(filename, opts); err != nil {
		return err
	}

	return mergo.Merge(o, opts, mergo.WithOverride)
}

func (o *Options) createManager() (*cluster.Manager, error) {
	infraProvisioner, err := provisioner.Create(o.Provisioner, &o.ProvisionerOptions)
	if err != nil {
		return nil, err
	}

	manifestRenderer, err := renderer.Create(o.Renderer, &o.RendererOptions)
	if err != nil {
		return nil, err
	}

	var credentialSource credentials.Source
	if !o.Credentials.Empty() {
		credentialSource = credentials.NewStaticSource(&o.Credentials)
	} else if outputter, ok := infraProvisioner.(provisioner.Outputter); ok {
		credentialSource = credentials.NewProvisionerOutputSource(outputter)
	} else {
		return nil, errors.New("please provide valid kubernetes credentials via the --cluster-* flags")
	}

	return cluster.NewManager(credentialSource, infraProvisioner, manifestRenderer), nil
}
