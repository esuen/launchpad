package store

import (
	"testing"

	"github.com/esuen/launchpad/internal/model"
)

func TestCreate(t *testing.T) {
	s := NewMemory()
	d, err := s.Create("api-server", "v1.2.3", "production")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if d.ServiceName != "api-server" {
		t.Errorf("expected service_name api-server, got %s", d.ServiceName)
	}
	if d.Version != "v1.2.3" {
		t.Errorf("expected version v1.2.3, got %s", d.Version)
	}
	if d.Environment != "production" {
		t.Errorf("expected environment production, got %s", d.Environment)
	}
	if d.Status != model.DeploymentStatusPending {
		t.Errorf("expected status pending, got %s", d.Status)
	}
	if d.ID == "" {
		t.Error("expected non-empty ID")
	}
	if d.CreatedAt.IsZero() {
		t.Error("expected non-zero CreatedAt")
	}
}

func TestGet(t *testing.T) {
	s := NewMemory()
	created, err := s.Create("api-server", "v1.0.0", "staging")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := s.Get(created.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != created.ID {
		t.Errorf("expected ID %s, got %s", created.ID, got.ID)
	}
}

func TestGetNotFound(t *testing.T) {
	s := NewMemory()

	_, err := s.Get("nonexistent")
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestList(t *testing.T) {
	s := NewMemory()
	s.Create("api-server", "v1.0.0", "production")
	s.Create("api-server", "v1.1.0", "staging")
	s.Create("web-app", "v2.0.0", "production")

	// List all.
	all, err := s.List("", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(all) != 3 {
		t.Errorf("expected 3 deployments, got %d", len(all))
	}

	// Filter by service.
	byService, _ := s.List("api-server", "")
	if len(byService) != 2 {
		t.Errorf("expected 2 deployments for api-server, got %d", len(byService))
	}

	// Filter by environment.
	byEnv, _ := s.List("", "production")
	if len(byEnv) != 2 {
		t.Errorf("expected 2 deployments for production, got %d", len(byEnv))
	}

	// Filter by both.
	byBoth, _ := s.List("api-server", "production")
	if len(byBoth) != 1 {
		t.Errorf("expected 1 deployment for api-server+production, got %d", len(byBoth))
	}
}
