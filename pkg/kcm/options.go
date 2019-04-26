package kcm

type Options struct {
	DryRun       bool   `json:"dryRun" yaml:"dryRun"`
	Manifest     string `json:"manifest" yaml:"manifest"`
	Values       string `json:"values" yaml:"values"`
	Deletions    string `json:"deletions" yaml:"deletions"`
	OnlyManifest bool   `json:"onlyManifest" yaml:"onlyManifest"`
}

type ProvisionerOptions struct {
	Terraform TerraformOptions `json:"terraform" yaml:"terraform"`
}

type TerraformOptions struct {
	Parallelism int `json:"parallelism" yaml:"parallelism"`
}

type RendererOptions struct {
	Helm HelmOptions `json:"helm" yaml:"helm"`
}

type HelmOptions struct {
	Chart string `json:"chart" yaml:"chart"`
}
