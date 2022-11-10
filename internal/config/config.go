package config

import (
	"go-herder/internal/api"
	"go-herder/internal/herder"
	"go-herder/internal/repository"
	"gopkg.in/yaml.v2"
	"io"
	"os"
)

type Config struct {
	DBConfig     repository.Config `yaml:"db"`
	HerderConfig herder.Config     `yaml:"herder"`
	APIConfig    api.Config        `yaml:"api"`
}

func New(configFile string) (*Config, error) {
	cfg := &Config{}
	file, err := os.Open(configFile)
	if err != nil {
		return cfg, err
	}
	cfgBytes, err := io.ReadAll(file)
	if err = file.Close(); err != nil {
		return cfg, err
	}
	err = yaml.Unmarshal(cfgBytes, &cfg)
	if cfg.DBConfig.DBFileName == "" {
		cfg.DBConfig.DBFileName = "database.db"
	}
	return cfg, err
}
