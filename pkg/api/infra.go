package api

// InfraOutput contains the output variables of an infrastructure manager.
type InfraOutput struct {
	Values map[string]interface{} `json:"values" yaml:"values"`
}

// NewInfraOutput creates a new InfraOutput value.
func NewInfraOutput() *InfraOutput {
	return &InfraOutput{Values: make(map[string]interface{})}
}
