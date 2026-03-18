package server

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/esuen/launchpad/internal/model"
	"github.com/esuen/launchpad/internal/store"
)

func newTestServer() *Server {
	logger := slog.Default()
	st := store.New()
	return New(logger, st)
}

func TestHealthz(t *testing.T) {
	srv := newTestServer()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()

	srv.Handler().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestReadyz(t *testing.T) {
	srv := newTestServer()
	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	w := httptest.NewRecorder()

	srv.Handler().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestCreateDeployment(t *testing.T) {
	srv := newTestServer()

	body := `{"service_name":"api-server","version":"v1.0.0","environment":"production"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/deployments/", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	srv.Handler().ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	var d model.Deployment
	if err := json.NewDecoder(w.Body).Decode(&d); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if d.ServiceName != "api-server" {
		t.Errorf("expected service_name api-server, got %s", d.ServiceName)
	}
	if d.Status != model.DeploymentStatusPending {
		t.Errorf("expected status pending, got %s", d.Status)
	}
}

func TestCreateDeploymentValidation(t *testing.T) {
	srv := newTestServer()

	body := `{"service_name":"api-server"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/deployments/", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	srv.Handler().ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestGetDeployment(t *testing.T) {
	srv := newTestServer()

	// Create one first.
	body := `{"service_name":"api-server","version":"v1.0.0","environment":"staging"}`
	createReq := httptest.NewRequest(http.MethodPost, "/api/v1/deployments/", bytes.NewBufferString(body))
	createReq.Header.Set("Content-Type", "application/json")
	createW := httptest.NewRecorder()
	srv.Handler().ServeHTTP(createW, createReq)

	var created model.Deployment
	json.NewDecoder(createW.Body).Decode(&created)

	// Get it.
	getReq := httptest.NewRequest(http.MethodGet, "/api/v1/deployments/"+created.ID, nil)
	getW := httptest.NewRecorder()
	srv.Handler().ServeHTTP(getW, getReq)

	if getW.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", getW.Code)
	}
}

func TestGetDeploymentNotFound(t *testing.T) {
	srv := newTestServer()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/deployments/nonexistent", nil)
	w := httptest.NewRecorder()
	srv.Handler().ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestListDeployments(t *testing.T) {
	srv := newTestServer()

	// Create two deployments.
	for _, body := range []string{
		`{"service_name":"api-server","version":"v1.0.0","environment":"production"}`,
		`{"service_name":"web-app","version":"v2.0.0","environment":"staging"}`,
	} {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/deployments/", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		srv.Handler().ServeHTTP(w, req)
	}

	// List all.
	req := httptest.NewRequest(http.MethodGet, "/api/v1/deployments/", nil)
	w := httptest.NewRecorder()
	srv.Handler().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var deployments []model.Deployment
	if err := json.NewDecoder(w.Body).Decode(&deployments); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(deployments) != 2 {
		t.Errorf("expected 2 deployments, got %d", len(deployments))
	}
}

func TestListDeploymentsEmpty(t *testing.T) {
	srv := newTestServer()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/deployments/", nil)
	w := httptest.NewRecorder()
	srv.Handler().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	// Should return empty array, not null.
	if w.Body.String() != "[]\n" {
		t.Errorf("expected empty array, got %s", w.Body.String())
	}
}
