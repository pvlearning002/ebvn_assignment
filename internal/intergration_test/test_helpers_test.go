package intergrationtest_test

import (
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/google/uuid"
	"github.com/pvlearning002/ebvn_assignment/internal/api"
	"github.com/pvlearning002/ebvn_assignment/internal/config"
	"github.com/pvlearning002/ebvn_assignment/pkg/utils"
)

func setupIntegrationAPI(t *testing.T) (*miniredis.Miniredis, api.Engine, *config.Config, *config.ConfigPaths, *config.ConfigHealthCheck) {
	t.Helper()

	redisServer := miniredis.RunT(t)
	cfg := config.CreateNewConfig()
	cfg.InstanceID = uuid.NewString()
	cfgStatus := config.CreateNewConfigStatus()
	cfgError := config.CreateNewConfigError()
	cfgPaths := config.CreateNewConfigPaths()
	cfgHealthCheck := config.CreateNewConfigHealthCheck()
	cfgParameter := config.CreateNewConfigParameter()
	cfgRedis := config.ConfigRedis{Address: redisServer.Addr()}

	logger := utils.NewLoggerUtils()
	response := utils.NewResponseUtils(*cfgError)
	stringUtils := utils.NewStringUtils(*cfg, *cfgError, logger)
	redisUtils := utils.NewRedisClientUtils(*cfg, cfgRedis, *cfgError)
	engine := api.NewEngine(
		cfg,
		cfgStatus,
		cfgError,
		cfgPaths,
		cfgHealthCheck,
		cfgParameter,
		&stringUtils,
		&redisUtils,
		&response,
		&logger,
	)

	return redisServer, engine, cfg, cfgPaths, cfgHealthCheck
}
