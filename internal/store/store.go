package store

import (
	"fmt"
	"sync"
	"time"

	"github.com/esuen/launchpad/internal/model"
	"github.com/google/uuid"
)

// Store is an in-memory deployment store.
type Store struct {
	mu          sync.RWMutex
	deployments map[string]model.Deployment
}

// New creates a new Store.
func New() *Store {
	return &Store{
		deployments: make(map[string]model.Deployment),
	}
}

// Create adds a new deployment and returns it.
func (s *Store) Create(serviceName, version, environment string) model.Deployment {
	d := model.Deployment{
		ID:          uuid.New().String(),
		ServiceName: serviceName,
		Version:     version,
		Environment: environment,
		Status:      model.DeploymentStatusPending,
		CreatedAt:   time.Now(),
	}

	s.mu.Lock()
	s.deployments[d.ID] = d
	s.mu.Unlock()

	return d
}

// ErrNotFound is returned when a deployment is not found.
var ErrNotFound = fmt.Errorf("deployment not found")

// Get returns a deployment by ID.
func (s *Store) Get(id string) (model.Deployment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	d, ok := s.deployments[id]
	if !ok {
		return model.Deployment{}, ErrNotFound
	}
	return d, nil
}

// List returns all deployments, optionally filtered by service name and/or environment.
func (s *Store) List(service, environment string) []model.Deployment {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var results []model.Deployment
	for _, d := range s.deployments {
		if service != "" && d.ServiceName != service {
			continue
		}
		if environment != "" && d.Environment != environment {
			continue
		}
		results = append(results, d)
	}
	return results
}
