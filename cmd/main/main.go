package main

import (
	apiPkg "go-herder/internal/api"
	configPkg "go-herder/internal/config"
	herderPkg "go-herder/internal/herder"
	"log"
)

func main() {
	cfg := configPkg.New()
	hrdr := herderPkg.New(cfg.HerderConfig)
	api := apiPkg.New(cfg.APIConfig, hrdr)
	if err := api.Run(); err != nil {
		log.Println("go-herder die with error:", err.Error())
	}
}
