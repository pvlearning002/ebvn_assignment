package main

import (
	"github.com/pvlearning002/ebvn_assignment/internal/api"
	"github.com/pvlearning002/ebvn_assignment/internal/config"
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
	engine := api.NewEngine(cfg, cfgStatus, cfgError, cfgPaths, cfgHealthCheck)
	if err := engine.Start(); err != nil {
		panic(err)
	}
}
