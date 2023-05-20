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
	Repos    *map[string]ConfigRepo `yaml:"repos"`
	Cloneurl string                 `yaml:"cloneurl"`
	Title    string                 `yaml:"title"`
	Seo      *map[string]string
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

	if config.Cloneurl == "" {
		config.Cloneurl = "http://localhost:8080"
	}
	if config.Title == "" {
		config.Title = "GetSrc"
	}
	if config.Seo == nil {
		config.Seo = &map[string]string{}
	}

	return config, nil
}
