package kubernetes

type ClusterOptions struct {
	Server     string `json:"server" yaml:"server"`
	Token      string `json:"token" yaml:"token"`
	Kubeconfig string `json:"kubeconfig" yaml:"kubeconfig"`
	Context    string `json:"context" yaml:"context"`
}

// Update tries to update the cluster config from values retrieved from the
// infrastructure manager. It will not overwrite config values that are already
// set.
func (o *ClusterOptions) Update(values map[string]interface{}) {
	if s, ok := values["server"].(string); ok && o.Server == "" {
		o.Server = s
	}

	if t, ok := values["token"].(string); ok && o.Token == "" {
		o.Token = t
	}

	if k, ok := values["kubeconfig"].(string); ok && o.Kubeconfig == "" {
		o.Kubeconfig = k
	}

	if v, ok := values["context"].(string); ok && o.Context == "" {
		o.Context = v
	}
}
