package store

import (
	"sync"
	"time"

	"github.com/esuen/launchpad/internal/model"
	"github.com/google/uuid"
)

// MemoryStore is an in-memory deployment store.
type MemoryStore struct {
	mu          sync.RWMutex
	deployments map[string]model.Deployment
}

// NewMemory creates a new MemoryStore.
func NewMemory() *MemoryStore {
	return &MemoryStore{
		deployments: make(map[string]model.Deployment),
	}
}

func (s *MemoryStore) Create(serviceName, version, environment string) (model.Deployment, error) {
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

	return d, nil
}

func (s *MemoryStore) Get(id string) (model.Deployment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	d, ok := s.deployments[id]
	if !ok {
		return model.Deployment{}, ErrNotFound
	}
	return d, nil
}

func (s *MemoryStore) List(service, environment string) ([]model.Deployment, error) {
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
	return results, nil
}
