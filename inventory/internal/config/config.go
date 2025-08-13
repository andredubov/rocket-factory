package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/andredubov/rocket-factory/inventory/internal/config/env"
)

var appConfig *config

type config struct {
	Logger     LoggerConfig
	GRPCServer GRPCConfig
	MongoDB    MongoDBConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	ufoGRPCCfg, err := env.NewGRPCConfig()
	if err != nil {
		return err
	}

	mongoDBCfg, err := env.NewMongoDBConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:     loggerCfg,
		GRPCServer: ufoGRPCCfg,
		MongoDB:    mongoDBCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
