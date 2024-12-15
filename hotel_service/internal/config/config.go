package config

import (
	"flag"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
	"log/slog"
)

type ServerConfig struct {
	Port   string `yaml:"port"`
	Prefix string `yaml:"prefix"`
}

func getConfigPath() (string, error) {
	slog.Info("Getting the config path")
	args := os.Args
	var configPath string

	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	flags.StringVar(&configPath, "config", "/hotel_service/internal/config/config.yaml", "path to configuration file")

	err := flags.Parse(args[1:])

	if err != nil {
		return "", err
	}

    slog.Info("Config path was successfully got")

	return configPath, nil
}

func GetServerConfig() (*ServerConfig, error) {
	slog.Info("Getting the server config")
	var cfg ServerConfig

	path, err := getConfigPath()

	if err != nil {
		return nil, err
	}

	cleanedPath := filepath.Clean(path)

	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	parentTop := filepath.Dir(wd)
	fullPath := parentTop + cleanedPath

	file, err := os.Open(fullPath)

	if err != nil {
		return nil, err
	}

	err = yaml.NewDecoder(file).Decode(&cfg)
	if err != nil {
		return nil, err
	}

	slog.Info("Server config was successfully got")

	return &cfg, nil
}
