package config

import (
	"fmt"
	"os"
)

type ServerConfig struct {
	Port          string
	Prefix        string
	EncryptionKey string
}

func GetServerConfig() (*ServerConfig, error) {
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		return nil, fmt.Errorf("environment variable SERVER_PORT is not set")
	}

	prefix := os.Getenv("SERVER_PREFIX")
	if prefix == "" {
		return nil, fmt.Errorf("environment variable SERVER_PREFIX is not set")
	}

	encryptionKey := os.Getenv("ENCRYPTION_KEY")
	if prefix == "" {
		return nil, fmt.Errorf("environment variable ENCRYPTION_KEY is not set")
	}

	return &ServerConfig{
		Port:          port,
		Prefix:        prefix,
		EncryptionKey: encryptionKey,
	}, nil
}
