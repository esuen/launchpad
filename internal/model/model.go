package model

import "time"

// DeploymentStatus represents the state of a deployment.
type DeploymentStatus string

const (
	DeploymentStatusPending   DeploymentStatus = "pending"
	DeploymentStatusRunning   DeploymentStatus = "running"
	DeploymentStatusSucceeded DeploymentStatus = "succeeded"
	DeploymentStatusFailed    DeploymentStatus = "failed"
)

// Deployment represents a single deployment of a service.
type Deployment struct {
	ID          string           `json:"id" db:"id"`
	ServiceName string           `json:"service_name" db:"service_name"`
	Version     string           `json:"version" db:"version"`
	Environment string           `json:"environment" db:"environment"`
	Status      DeploymentStatus `json:"status" db:"status"`
	CreatedAt   time.Time        `json:"created_at" db:"created_at"`
}
