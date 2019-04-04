package config

type Config struct {
	DryRun       bool   `json:"dryRun"`
	OnlyManifest bool   `json:"onlyManifest"`
	WorkingDir   string `json:"workingDir"`
	Manifest     string `json:"manifest"`
	Kubeconfig   string `json:"kubeconfig"`

	Deletions string `json:"deletions"`

	Terraform TerraformConfig `json:"terraform"`

	Helm HelmConfig `json:"helm"`
}

type TerraformConfig struct {
	AutoApprove      bool `json"autoApprove"`
	Parallelism      int  `json:"parallelism"`
	DetailedExitCode bool `json:"detailedExitCode"`
}

type HelmConfig struct {
	Values string `json:"values"`
	Chart  string `json:"chart"`
}

func (c *Config) ApplyDefaults() {
	if c.Manifest == "" {
		c.Manifest = c.WorkingDir + "/manifest.yaml"
	}

	if c.Deletions == "" {
		c.Deletions = c.WorkingDir + "/deletions.yaml"
	}

	if c.Helm.Values == "" {
		c.Helm.Values = c.WorkingDir + "/values.yaml"
	}

	if c.Helm.Chart == "" {
		c.Helm.Chart = c.WorkingDir + "/cluster"
	}
}
