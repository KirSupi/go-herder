package config

import (
	"go-herder/internal/api"
	"go-herder/internal/herder"
)

type Config struct {
	APIConfig    api.Config
	HerderConfig herder.Config
}

func New() *Config {
	return &Config{}
}
