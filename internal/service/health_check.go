package service

import (
	"github.com/google/uuid"
	"github.com/pvlearning002/ebvn_assignment/internal/config"
	"github.com/pvlearning002/ebvn_assignment/internal/model"
)

//go:generate mockery --name=HealthCheck --filename=health_check_mck.go
type HealthCheck interface {
	GenerateHealthCheck() (*model.HealthCheckOutput, error)
}

type healthCheckService struct {
	cfg            config.Config
	cfgStatus      config.ConfigStatus
	cfgHealthCheck config.ConfigHealthCheck
}

func NewHealthCheck(cfg config.Config, cfgStatus config.ConfigStatus, cfgHealthCheck config.ConfigHealthCheck) HealthCheck {
	return &healthCheckService{
		cfg:            cfg,
		cfgStatus:      cfgStatus,
		cfgHealthCheck: cfgHealthCheck,
	}
}

func (self *healthCheckService) parseUUIDfromString(uuidValue string) (uuid.UUID, error) {
	result, err := uuid.Parse(uuidValue)
	if err != nil {
		return uuid.Nil, err
	}
	return result, nil
}

func (self *healthCheckService) GenerateHealthCheck() (*model.HealthCheckOutput, error) {
	instanceID, err := self.parseUUIDfromString(self.cfg.InstanceID)
	if err != nil {
		return nil, err
	}

	result := &model.HealthCheckOutput{
		Message:     self.cfgStatus.StatusOK,
		ServiceName: self.cfgHealthCheck.BookmarkService,
		InstanceID:  instanceID,
	}
	return result, nil
}
