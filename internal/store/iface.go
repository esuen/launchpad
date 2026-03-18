package store

import (
	"fmt"

	"github.com/esuen/launchpad/internal/model"
)

// ErrNotFound is returned when a deployment is not found.
var ErrNotFound = fmt.Errorf("deployment not found")

// Store defines the interface for deployment storage.
type Store interface {
	Create(serviceName, version, environment string) (model.Deployment, error)
	Get(id string) (model.Deployment, error)
	List(service, environment string) ([]model.Deployment, error)
}
