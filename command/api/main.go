package main

import (
	"github.com/pvlearning002/ebvn_assignment/internal/api"
	"github.com/pvlearning002/ebvn_assignment/internal/config"
)

func main() {
	cfg := config.CreateNewConfig()
	engine := api.NewEngine(cfg)
	if err := engine.Start(); err != nil {
		panic(err)
	}
}
