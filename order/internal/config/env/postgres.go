package env

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type postgresEnvConfig struct {
	Host         string `env:"PG_HOST,required"`
	Port         string `env:"PG_PORT,required"`
	User         string `env:"PG_USER,required"`
	DatabaseName string `env:"PG_DB,required"`
	Password     string `env:"PG_PASSWORD,required"`
	SSLMode      string `env:"PG_SSL_MODE,required"`
	MigrationDir string `env:"MIGRATION_DIR,required"`
}

type postgresConfig struct {
	raw postgresEnvConfig
}

func NewPostgresDBConfig() (*postgresConfig, error) {
	var raw postgresEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &postgresConfig{raw: raw}, nil
}

func (cfg *postgresConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		cfg.raw.Host,
		cfg.raw.Port,
		cfg.raw.DatabaseName,
		cfg.raw.User,
		cfg.raw.Password,
		cfg.raw.SSLMode,
	)
}

func (cfg *postgresConfig) MigrationDirectory() string {
	return cfg.raw.MigrationDir
}
