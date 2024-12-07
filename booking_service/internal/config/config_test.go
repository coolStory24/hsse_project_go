package config

//
//import (
//	"flag"
//	"os"
//	"path/filepath"
//	"testing"
//
//	"github.com/stretchr/testify/assert"
//)
//
//func createTempConfigFile(t *testing.T, content string) string {
//	t.Helper()
//	tmpFile, err := os.CreateTemp("", "test-config-*.yaml")
//	if err != nil {
//		t.Fatalf("Failed to create temporary file: %v", err)
//	}
//	defer tmpFile.Close()
//
//	_, err = tmpFile.WriteString(content)
//	if err != nil {
//		t.Fatalf("Failed to write to temporary file: %v", err)
//	}
//
//	return tmpFile.Name()
//}
//
//func TestGetServerConfig_ValidConfig(t *testing.T) {
//	content := `
//port: "8080"
//prefix: "/api"
//`
//	configPath := createTempConfigFile(t, content)
//	defer os.Remove(configPath)
//
//	// Simulate the `--config` flag
//	os.Args = []string{"test", "--config", configPath}
//	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
//
//	cfg, err := GetServerConfig()
//	assert.NoError(t, err)
//	assert.Equal(t, "8080", cfg.Port)
//	assert.Equal(t, "/api", cfg.Prefix)
//}
//
//func TestGetServerConfig_MissingFile(t *testing.T) {
//	// Simulate a missing file
//	os.Args = []string{"test", "--config", "non-existent.yaml"}
//	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
//
//	_, err := GetServerConfig()
//	assert.Error(t, err)
//}
//
//func TestGetServerConfig_InvalidYAML(t *testing.T) {
//	content := `
//port: "8080"
//prefix: "/api
//`
//	configPath := createTempConfigFile(t, content)
//	defer os.Remove(configPath)
//
//	// Simulate the `--config` flag
//	os.Args = []string{"test", "--config", configPath}
//	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
//
//	_, err := GetServerConfig()
//	assert.Error(t, err)
//}
//
//func TestGetServerConfig_DefaultPath(t *testing.T) {
//	// Ensure the default path doesn't exist
//	defaultPath := filepath.Clean("/booking_service/internal/config/config.yaml")
//	if _, err := os.Stat(defaultPath); err == nil {
//		t.Skip("Default path exists, skipping test")
//	}
//
//	os.Args = []string{"test"}
//	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
//
//	_, err := GetServerConfig()
//	assert.Error(t, err)
//}
