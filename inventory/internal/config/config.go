package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/andredubov/rocket-factory/inventory/internal/config/env"
	iamClient "github.com/andredubov/rocket-factory/inventory/internal/config/env/iam"
)

var appConfig *config

type config struct {
	Logger     LoggerConfig
	GRPCServer GRPCConfig
	MongoDB    MongoDBConfig
	IAMClient  GRPCConfig
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

	grpcServerConfig, err := env.NewGRPCConfig()
	if err != nil {
		return err
	}

	iamClientConfig, err := iamClient.NewGRPCConfig()
	if err != nil {
		return err
	}

	mongoDBConfig, err := env.NewMongoDBConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:     loggerConfig,
		GRPCServer: grpcServerConfig,
		MongoDB:    mongoDBConfig,
		IAMClient:  iamClientConfig,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
