package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/andredubov/rocket-factory/payment/internal/config/env"
	iamClient "github.com/andredubov/rocket-factory/payment/internal/config/env/iam"
)

var appConfig *config

type config struct {
	Logger     LoggerConfig
	GRPCServer GRPCConfig
	IAMClient  GRPCConfig
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

	grpcServerCfg, err := env.NewGRPCConfig()
	if err != nil {
		return err
	}

	iamClientConfig, err := iamClient.NewGRPCConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:     loggerCfg,
		GRPCServer: grpcServerCfg,
		IAMClient:  iamClientConfig,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
