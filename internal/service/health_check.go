package service

import (
	"github.com/google/uuid"
	"github.com/pvlearning002/ebvn_assignment/internal/config"
	"github.com/pvlearning002/ebvn_assignment/internal/model"
)

//go:generate mockery --name=HealthCheck --filename=health_check_mck.go
type HealthCheck interface {
	GenerateHealthCheck() (*model.HealthCheck, error)
}

type healthCheckService struct {
	cfg config.Config
}

func NewHealthCheck(cfg config.Config) HealthCheck {
	return &healthCheckService{
		cfg: cfg,
	}
}

func (self *healthCheckService) GenerateHealthCheck() (*model.HealthCheck, error) {
	result := &model.HealthCheck{
		Message:     self.cfg.StatusOK,
		ServiceName: self.cfg.BookmarkService,
		InstanceID:  uuid.New(),
	}
	return result, nil
}
