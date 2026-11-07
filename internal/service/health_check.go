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
	cfgStatus      config.ConfigStatus
	cfgHealthCheck config.ConfigHealthCheck
}

func NewHealthCheck(cfgStatus config.ConfigStatus, cfgHealthCheck config.ConfigHealthCheck) HealthCheck {
	return &healthCheckService{
		cfgStatus:      cfgStatus,
		cfgHealthCheck: cfgHealthCheck,
	}
}

func (self *healthCheckService) GenerateHealthCheck() (*model.HealthCheck, error) {
	result := &model.HealthCheck{
		Message:     self.cfgStatus.StatusOK,
		ServiceName: self.cfgHealthCheck.BookmarkService,
		InstanceID:  uuid.New(),
	}
	return result, nil
}
