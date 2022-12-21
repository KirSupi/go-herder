package herder

import (
	"gopkg.in/yaml.v2"
	"io"
	"os"
)

type Config struct {
	HerderConfig HerderConfig `yaml:"herder"`
	APIConfig    APIConfig    `yaml:"api"`
}

func NewConfigFromFile(configFilePath string) (*Config, error) {
	cfg := &Config{}
	file, err := os.Open(configFilePath)
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
