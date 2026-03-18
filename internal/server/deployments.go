package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/esuen/launchpad/internal/model"
	"github.com/esuen/launchpad/internal/store"
	"github.com/go-chi/chi/v5"
)

type createDeploymentRequest struct {
	ServiceName string `json:"service_name"`
	Version     string `json:"version"`
	Environment string `json:"environment"`
}

func (s *Server) handleCreateDeployment(w http.ResponseWriter, r *http.Request) {
	var req createDeploymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	if req.ServiceName == "" || req.Version == "" || req.Environment == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "service_name, version, and environment are required"})
		return
	}

	d := s.store.Create(req.ServiceName, req.Version, req.Environment)
	writeJSON(w, http.StatusCreated, d)
}

func (s *Server) handleGetDeployment(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	d, err := s.store.Get(id)
	if errors.Is(err, store.ErrNotFound) {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "deployment not found"})
		return
	}

	writeJSON(w, http.StatusOK, d)
}

func (s *Server) handleListDeployments(w http.ResponseWriter, r *http.Request) {
	service := r.URL.Query().Get("service")
	environment := r.URL.Query().Get("environment")

	deployments := s.store.List(service, environment)
	if deployments == nil {
		deployments = []model.Deployment{}
	}
	writeJSON(w, http.StatusOK, deployments)
}
