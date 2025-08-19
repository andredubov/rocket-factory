package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/andredubov/rocket-factory/order/internal/config/env"
	inventoryClient "github.com/andredubov/rocket-factory/order/internal/config/env/inventory"
	paymentClient "github.com/andredubov/rocket-factory/order/internal/config/env/payment"
)

var appConfig *config

type config struct {
	Logger          LoggerConfig
	HTTPServer      HTTPConfig
	PostgresDB      PostgresDBConfig
	InventoryClient GRPCConfig
	PaymentClient   GRPCConfig
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

	httpConfig, err := env.NewHTTPConfig()
	if err != nil {
		return err
	}

	postgresDBConfig, err := env.NewPostgresDBConfig()
	if err != nil {
		return err
	}

	invetoryConfig, err := inventoryClient.NewGRPCConfig()
	if err != nil {
		return err
	}

	paymentConfig, err := paymentClient.NewGRPCConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:          loggerConfig,
		HTTPServer:      httpConfig,
		PostgresDB:      postgresDBConfig,
		InventoryClient: invetoryConfig,
		PaymentClient:   paymentConfig,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
