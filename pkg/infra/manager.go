package infra

import "github.com/martinohmann/cluster-manager/pkg/api"

type Manager interface {
	Apply() (*api.InfraOutput, error)
}
