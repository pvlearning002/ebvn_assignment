package main

import (
	"fmt"

	"github.com/pvlearning002/ebvn_assignment/internal/config"
	"github.com/pvlearning002/ebvn_assignment/internal/service"
)

func main() {
	cfg := config.CreateNewConfig()
	fmt.Println("App Port:", cfg.AppPort)
	genPassService := service.NewGenPass(*cfg)
	generatedPass, err := genPassService.GeneratePassword()
	if err != nil {
		fmt.Println("Error generating password:", err)
		return
	}
	fmt.Println("Generated Password:", generatedPass)

}
