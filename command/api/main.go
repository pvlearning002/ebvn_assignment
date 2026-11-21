package main

// cspell:ignore pvlearning ebvn
import (
	"uuid"

	"github.com/pvlearning002/ebvn_assignment/internal/api"
	"github.com/pvlearning002/ebvn_assignment/internal/config"
	"github.com/pvlearning002/ebvn_assignment/pkg/utils/utils_log"
	redis_utils "github.com/pvlearning002/ebvn_assignment/pkg/utils/utils_redis"
)

// @title Bookmark Management API
// @version 1.0.0
// @description This is a sample server for managing bookmarks.
// @BasePath /

func main() {
	cfg := config.CreateNewConfig()
	cfgStatus := config.CreateNewConfigStatus()
	cfgError := config.CreateNewConfigError()
	cfgPaths := config.CreateNewConfigPaths()
	cfgHealthCheck := config.CreateNewConfigHealthCheck()
	cfgParameter := config.CreateNewConfigParameter()
	cfgRedis := config.CreateNewConfigRedis()
	custom_logger := utils_log.NewCustomLogger()

	if cfg.InstanceID == cfg.StringEmpty {
		cfg.InstanceID = uuid.New().String()
	}

	custom_logger.SetLogLevel(cfg.LogLevel)

	redisClient, err := redis_utils.NewRedisClient(*cfgRedis)
	if err != nil {
		panic(err)
	}
	defer redisClient.Close()

	engine := api.NewEngine(cfg, cfgStatus, cfgError, cfgPaths, cfgHealthCheck, cfgParameter, redisClient, &custom_logger)
	if err := engine.Start(); err != nil {
		panic(err)
	}
}
