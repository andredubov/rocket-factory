package config

import (
	"errors"
	"flag"
	"time"

	"github.com/joho/godotenv"
)

const (
	// ConfigPathFlagName is cli flag for config file path
	ConfigPathFlagName = "config-path"
	// ConfigPathFlagValue is cli flag default value for config path
	ConfigPathFlagValue = ""
)

var (
	// ErrEmptyConfigFilePath is error config path is empty
	ErrEmptyConfigFilePath = errors.New("config path is empty")
	// ErrConfigFileDoesNotExist is error config file does't exist
	ErrConfigFileDoesNotExist = errors.New("config file does't exist")
)

// GRPCConfig interface
type GRPCConfig interface {
	Address() string
}

// HTTPConfig interface
type HTTPConfig interface {
	Address() string
	ReadHeaderTimeout() time.Duration
}

// PostgresConfig interface
type PostgresConfig interface {
	DSN() string
	MigrationDirectory() string
}

// Load enviriment variables from *.env file
func Load() error {
	var configPath string

	if flag.Lookup(ConfigPathFlagName) == nil {
		flag.StringVar(&configPath, ConfigPathFlagName, ConfigPathFlagValue, "path to env file")
	}
	flag.Parse()

	if configPath == ConfigPathFlagValue {
		return nil
	}

	err := godotenv.Load(configPath)
	if err != nil {
		return err
	}

	return nil
}
