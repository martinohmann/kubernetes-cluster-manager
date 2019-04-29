package kcm

// Options are used to configure the cluster manager.
type Options struct {
	DryRun        bool   `json:"dryRun" yaml:"dryRun"`
	Values        string `json:"values" yaml:"values"`
	Deletions     string `json:"deletions" yaml:"deletions"`
	ManifestsDir  string `json:"manifestsDir" yaml:"manifestsDir"`
	SkipManifests bool   `json:"skipManifests" yaml:"skipManifests"`
	OnlyChanges   bool   `json:"onlyChanges" yaml:"onlyChanges"`
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
	ChartsDir string `json:"chartsDir" yaml:"chartsDir"`
}
