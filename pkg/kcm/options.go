package kcm

// Options are used to configure the cluster manager.
type Options struct {
	DryRun        bool   `json:"dryRun" yaml:"dryRun"`
	Manifest      string `json:"manifest" yaml:"manifest"`
	Values        string `json:"values" yaml:"values"`
	Deletions     string `json:"deletions" yaml:"deletions"`
	SkipManifests bool   `json:"skipManifests" yaml:"skipManifests"`
}

// ProvisionerOptions are made available to infrastructure provisioners.
type ProvisionerOptions struct {
	Terraform TerraformOptions `json:"terraform" yaml:"terraform"`
}

// TerraformOptions configure the terraform infrastructure provisioner.
type TerraformOptions struct {
	Parallelism int `json:"parallelism" yaml:"parallelism"`
}

// RendererOptions are made available to manifest renderers.
type RendererOptions struct {
	Helm HelmOptions `json:"helm" yaml:"helm"`
}

// HelmOptions configure the helm manifest renderer.
type HelmOptions struct {
	Chart string `json:"chart" yaml:"chart"`
}
