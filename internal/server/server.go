package server

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/esuen/launchpad/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Server is the main application server.
type Server struct {
	logger *slog.Logger
	store  store.Store
}

// New creates a new Server.
func New(logger *slog.Logger, store store.Store) *Server {
	return &Server{logger: logger, store: store}
}

// Handler returns the HTTP handler for the server.
func (s *Server) Handler() http.Handler {
	r := chi.NewRouter()

	r.Use(loggingMiddleware)
	r.Use(metricsMiddleware)

	r.Get("/healthz", s.handleHealthz)
	r.Get("/readyz", s.handleReadyz)
	r.Handle("/metrics", promhttp.Handler())

	r.Route("/api/v1/deployments", func(r chi.Router) {
		r.Post("/", s.handleCreateDeployment)
		r.Get("/", s.handleListDeployments)
		r.Get("/{id}", s.handleGetDeployment)
	})

	return r
}

func (s *Server) handleHealthz(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleReadyz(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("failed to write response", "error", err)
	}
}
