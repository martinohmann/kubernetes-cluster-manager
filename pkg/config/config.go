package config

import "github.com/imdario/mergo"

// Config holds the configuration for kcm.
type Config struct {
	Debug            bool   `json:"debug" yaml:"debug"`
	DryRun           bool   `json:"dryRun" yaml:"dryRun"`
	OnlyManifest     bool   `json:"onlyManifest" yaml:"onlyManifest"`
	WorkingDir       string `json:"workingDir" yaml:"workingDir"`
	Manifest         string `json:"manifest" yaml:"manifest"`
	Values           string `json:"values" yaml:"values"`
	Deletions        string `json:"deletions" yaml:"deletions"`
	ManifestRenderer string `json:"manifestRenderer" yaml:"maniestRenderer"`
	InfraManager     string `json:"infraManager" yaml:"infraManager"`

	Cluster   ClusterConfig   `json:"cluster" yaml:"cluster"`
	Terraform TerraformConfig `json:"terraform" yaml:"terraform"`
	Helm      HelmConfig      `json:"helm" yaml:"helm"`
}

// Merge merges other into c. Fields in c that do not have their default value
// will not be overwritten.
func (c *Config) Merge(other *Config) error {
	return mergo.Merge(c, other)
}

// ApplyDefaults applies sane default values to fields that are not explicitly
// set.
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

// ClusterConfig holds configuration accessing a kubernetes cluster.
type ClusterConfig struct {
	Server     string `json:"server" yaml:"server"`
	Token      string `json:"token" yaml:"token"`
	Kubeconfig string `json:"kubeconfig" yaml:"kubeconfig"`
	Context    string `json:"context" yaml:"context"`
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

// TerraformConfig holds flags for terraform commands.
type TerraformConfig struct {
	Parallelism int `json:"parallelism" yaml:"parallelism"`
}

// HelmConfig holds the chart to be used by helm.
type HelmConfig struct {
	Chart string `json:"chart" yaml:"chart"`
}
