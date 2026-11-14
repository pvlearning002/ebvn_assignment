package main

// cspell:ignore pvlearning ebvn
import (
	"github.com/pvlearning002/ebvn_assignment/internal/api"
	"github.com/pvlearning002/ebvn_assignment/internal/config"
	"github.com/pvlearning002/ebvn_assignment/pkg/logger"
	"github.com/pvlearning002/ebvn_assignment/pkg/redis_utils"
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
	logger.SetLogLevel(cfg.LogLevel)
	redisClient, err := redis_utils.NewRedisClient(cfg.StringEmpty)
	if err != nil {
		panic(err)
	}
	defer redisClient.Close()

	engine := api.NewEngine(cfg, cfgStatus, cfgError, cfgPaths, cfgHealthCheck, cfgParameter, redisClient)
	if err := engine.Start(); err != nil {
		panic(err)
	}
}
