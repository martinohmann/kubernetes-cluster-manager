package infra

import (
	"github.com/martinohmann/kubernetes-cluster-manager/infra/terraform"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/api"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/command"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/config"
	"github.com/pkg/errors"
)

// Manager is the interface for a cloud infrastructure manager.
type Manager interface {
	// Apply will apply changes to the infrastructure. It will automatically
	// create or update a kubernetes cluster.
	Apply() error

	// Plan will plan changes to the infrastructure without acutally applying
	// them..
	Plan() error

	// GetValues obtains output values from the infrastructure manager.
	// These values are made available during kubernetes manifest
	// renderering.
	GetValues() (api.Values, error)

	// Destroy performs all actions needed to destroy a kubernetes cluster.
	Destroy() error
}

// CreateManager creates an infrastructure manager based on the config.
func CreateManager(cfg *config.Config, executor command.Executor) (Manager, error) {
	switch cfg.InfraManager {
	case "terraform":
		return terraform.NewInfraManager(&cfg.Terraform, executor), nil
	default:
		return nil, errors.Errorf("unsupported infrastructure manager: %s", cfg.InfraManager)
	}
}
