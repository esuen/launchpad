package store

import (
	"fmt"

	"github.com/esuen/launchpad/internal/model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// PostgresStore is a Postgres-backed deployment store.
type PostgresStore struct {
	db *sqlx.DB
}

// NewPostgres creates a new PostgresStore.
func NewPostgres(db *sqlx.DB) *PostgresStore {
	return &PostgresStore{db: db}
}

func (s *PostgresStore) Create(serviceName, version, environment string) (model.Deployment, error) {
	var d model.Deployment
	err := s.db.QueryRowx(
		`INSERT INTO deployments (id, service_name, version, environment, status)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, service_name, version, environment, status, created_at`,
		uuid.New().String(), serviceName, version, environment, model.DeploymentStatusPending,
	).StructScan(&d)
	if err != nil {
		return model.Deployment{}, err
	}
	return d, nil
}

func (s *PostgresStore) Get(id string) (model.Deployment, error) {
	var d model.Deployment
	err := s.db.Get(&d, `SELECT * FROM deployments WHERE id = $1`, id)
	if err != nil {
		return model.Deployment{}, ErrNotFound
	}
	return d, nil
}

func (s *PostgresStore) List(service, environment string) ([]model.Deployment, error) {
	query := `SELECT * FROM deployments WHERE 1=1`
	args := []any{}
	argN := 1

	if service != "" {
		query += fmt.Sprintf(` AND service_name = $%d`, argN)
		args = append(args, service)
		argN++
	}
	if environment != "" {
		query += fmt.Sprintf(` AND environment = $%d`, argN)
		args = append(args, environment)
		argN++
	}

	query += ` ORDER BY created_at DESC`

	var deployments []model.Deployment
	if err := s.db.Select(&deployments, query, args...); err != nil {
		return nil, err
	}
	return deployments, nil
}
