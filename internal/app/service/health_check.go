package service

import (
	"github.com/google/uuid"
	"github.com/pvlearning002/ebvn_assignment/internal/app/model"
	"github.com/pvlearning002/ebvn_assignment/internal/config"
	"github.com/pvlearning002/ebvn_assignment/pkg/utils"
)

//go:generate mockery --name=HealthCheck --filename=health_check_mck.go
type HealthCheck interface {
	GenerateHealthCheck() (*model.HealthCheckOutput, error)
}

type healthCheckService struct {
	cfg            config.Config
	cfgStatus      config.ConfigStatus
	cfgHealthCheck config.ConfigHealthCheck
	redisUtils     utils.RedisClientUtils
}

func NewHealthCheck(cfg config.Config,
	cfgStatus config.ConfigStatus,
	cfgHealthCheck config.ConfigHealthCheck,
	redisUtils utils.RedisClientUtils) HealthCheck {
	return &healthCheckService{
		cfg:            cfg,
		cfgStatus:      cfgStatus,
		cfgHealthCheck: cfgHealthCheck,
		redisUtils:     redisUtils,
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
	isRedisConnect := self.redisUtils.IsRedisContainerRunning()

	result := &model.HealthCheckOutput{
		Message:          self.cfgStatus.StatusOK,
		ServiceName:      self.cfgHealthCheck.BookmarkService,
		InstanceID:       instanceID,
		IsRedisConnected: isRedisConnect,
	}
	return result, nil
}
