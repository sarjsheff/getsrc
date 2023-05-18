package getsrc

import (
	"os"

	"gopkg.in/yaml.v3"
)

type ConfigRepo struct {
	Path        string `yaml:"path"`
	Description string `yaml:"description"`
}

type Config struct {
	Repos *map[string]ConfigRepo `yaml:"repos"`
}

func NewConfig(configPath string) (*Config, error) {
	config := &Config{}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)

	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}
