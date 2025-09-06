package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/andredubov/rocket-factory/iam/internal/config/env"
)

var appConfig *config

type config struct {
	Logger     LoggerConfig
	GRPCServer GRPCConfig
	PostgresDB PostgresDBConfig
	RedisCache RedisConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	loggerConfig, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	grpcConfig, err := env.NewGRPCConfig()
	if err != nil {
		return err
	}

	postgresCfg, err := env.NewPostgresDBConfig()
	if err != nil {
		return err
	}
	redisCfg, err := env.NewRedisConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:     loggerConfig,
		GRPCServer: grpcConfig,
		PostgresDB: postgresCfg,
		RedisCache: redisCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
