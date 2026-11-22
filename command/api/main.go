package main

// cspell:ignore pvlearning ebvn
import (
	"github.com/google/uuid"

	"github.com/pvlearning002/ebvn_assignment/internal/api"
	"github.com/pvlearning002/ebvn_assignment/internal/config"
	"github.com/pvlearning002/ebvn_assignment/pkg/utils"
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
	logger := utils.NewLoggerUtils()
	stringUtils := utils.NewStringUtils(*cfg, *cfgError, logger)
	redisUtils := utils.NewRedisClientUtils(*cfg, *cfgRedis, *cfgError)

	if cfg.InstanceID == cfg.StringEmpty {
		cfg.InstanceID = uuid.New().String()
	}

	logger.SetLogLevel(cfg.LogLevel)

	engine := api.NewEngine(cfg, cfgStatus, cfgError, cfgPaths, cfgHealthCheck, cfgParameter, &stringUtils, &redisUtils, &logger)
	if err := engine.Start(); err != nil {
		panic(err)
	}
}
