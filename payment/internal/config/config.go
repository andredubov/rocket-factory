package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/andredubov/rocket-factory/payment/internal/config/env"
)

var appConfig *config

type config struct {
	Logger     LoggerConfig
	GRPCServer GRPCConfig
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

	appConfig = &config{
		Logger:     loggerCfg,
		GRPCServer: ufoGRPCCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
