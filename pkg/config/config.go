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

type ClusterConfig struct {
	Server     string `json:"server"`
	Token      string `json:"token"`
	Kubeconfig string `json:"kubeconfig"`
	Context    string `json:"context"`
}

// Update tries to update the cluster config from values retrieved from the
// infrastructure manager. It will not overwrite config values that are already
// set.
func (c *ClusterConfig) Update(values map[string]interface{}) {
	if s, ok := values["server"].(string); ok && c.Server == "" {
		c.Server = s
	}

	if t, ok := values["token"].(string); ok && c.Token == "" {
		c.Token = t
	}

	if k, ok := values["kubeconfig"].(string); ok && c.Kubeconfig == "" {
		c.Kubeconfig = k
	}

	if v, ok := values["context"].(string); ok && c.Context == "" {
		c.Context = v
	}
}

type TerraformConfig struct {
	Parallelism int `json:"parallelism"`
}

type HelmConfig struct {
	Chart string `json:"chart"`
}
