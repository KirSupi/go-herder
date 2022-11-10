package main

import (
	apiPkg "go-herder/internal/api"
	configPkg "go-herder/internal/config"
	herderPkg "go-herder/internal/herder"
	"go-herder/internal/repository"
	"log"
	"os"
)

func main() {
	if len(os.Args) > 2 {
		log.Fatalln("too many args")
	}
	configFileArg := "herder-config.yml"
	if len(os.Args) == 2 {
		configFileArg = os.Args[1]
	}
	cfg, err := configPkg.New(configFileArg)
	if err != nil {
		log.Fatalln("error when loading config:", err.Error())
	}
	repo, err := repository.New(cfg.DBConfig)
	if err != nil {
		log.Fatalln("error when creating repository:", err.Error())
	}
	hrdr := herderPkg.New(cfg.HerderConfig, repo)
	if err = hrdr.Init(); err != nil {
		log.Fatalln("error on init Herder:", err.Error())
	}
	//if err = hrdr.RunAll(); err != nil {
	//	log.Fatalln("error on run Herder:", err.Error())
	//}
	api := apiPkg.New(cfg.APIConfig, hrdr, repo)
	if err = api.Run(); err != nil {
		log.Println("go-herder die with error:", err.Error())
	}
}
