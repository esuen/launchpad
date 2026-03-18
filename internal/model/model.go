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
	ID          string           `json:"id"`
	ServiceName string           `json:"service_name"`
	Version     string           `json:"version"`
	Environment string           `json:"environment"`
	Status      DeploymentStatus `json:"status"`
	CreatedAt   time.Time        `json:"created_at"`
}
