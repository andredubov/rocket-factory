package config

import "time"

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type GRPCConfig interface {
	Address() string
}

type HTTPConfig interface {
	Address() string
	ReadHeaderTimeout() time.Duration
}

type PostgresDBConfig interface {
	DSN() string
	MigrationDirectory() string
}
