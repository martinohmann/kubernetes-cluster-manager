package config

type Config struct {
	Debug            bool   `json:"debug"`
	DryRun           bool   `json:"dryRun"`
	OnlyManifest     bool   `json:"onlyManifest"`
	WorkingDir       string `json:"workingDir"`
	Manifest         string `json:"manifest"`
	Values           string `json:"values"`
	Deletions        string `json:"deletions"`
	ManifestRenderer string `json:"manifestRenderer"`
	InfraManager     string `json:"infraManager"`

	Cluster   ClusterConfig   `json:"cluster"`
	Terraform TerraformConfig `json:"terraform"`
	Helm      HelmConfig      `json:"helm"`
}

type ClusterConfig struct {
	Server     string `json:"server"`
	Token      string `json:"token"`
	Kubeconfig string `json:"kubeconfig"`
}

type TerraformConfig struct {
	Parallelism int `json:"parallelism"`
}

type HelmConfig struct {
	Chart string `json:"chart"`
}

func (c *Config) ApplyDefaults() {
	if c.Manifest == "" {
		c.Manifest = c.WorkingDir + "/manifest.yaml"
	}

	if c.Deletions == "" {
		c.Deletions = c.WorkingDir + "/deletions.yaml"
	}

	if c.Values == "" {
		c.Values = c.WorkingDir + "/values.yaml"
	}

	if c.Helm.Chart == "" {
		c.Helm.Chart = c.WorkingDir + "/cluster"
	}
}
