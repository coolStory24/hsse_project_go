package config

import (
	"flag"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type ServerConfig struct {
	Port   string `yaml:"port"`
	Prefix string `yaml:"prefix"`
}

func getConfigPath() (string, error) {
	args := os.Args
	var configPath string

	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	flags.StringVar(&configPath, "config", "../internal/config/config.yaml", "path to configuration file")

	err := flags.Parse(args[1:])

	if err != nil {
		return "", err
	}

	return configPath, nil
}

func GetServerConfig() (*ServerConfig, error) {
	var cfg ServerConfig

	path, err := getConfigPath()

	if err != nil {
		return nil, err
	}

	cleanedPath := filepath.Clean(path)

	file, err := os.Open(cleanedPath)

	if err != nil {
		return nil, err
	}

	err = yaml.NewDecoder(file).Decode(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
