package config

import (
	"go-herder/internal/api"
	"go-herder/internal/herder"
	"gopkg.in/yaml.v2"
	"io"
	"os"
)

type Config struct {
	HerderConfig herder.Config `yaml:"herder"`
	APIConfig    api.Config    `yaml:"api"`
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
	return cfg, err
}
