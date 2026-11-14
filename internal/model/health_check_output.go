package model

import (
	"github.com/google/uuid"
)

type HealthCheckOutput struct {
	Message     string    `json:"message"`
	ServiceName string    `json:"service_name"`
	InstanceID  uuid.UUID `json:"instance_id"`
}
