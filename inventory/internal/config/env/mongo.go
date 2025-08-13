package env

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type mongoEnvConfig struct {
	Username    string `env:"MONGO_INITDB_ROOT_USERNAME,required"`
	Password    string `env:"MONGO_INITDB_ROOT_PASSWORD,required"`
	Host        string `env:"MONGO_HOST,required"`
	Port        string `env:"MONGO_PORT,required"`
	AuthSource  string `env:"MONGO_AUTH_SOURCE,required"`
	MongoDBName string `env:"MONGO_INITDB_DATABASE,required"`
}

type mongoConfig struct {
	raw mongoEnvConfig
}

func NewMongoDBConfig() (*mongoConfig, error) {
	var raw mongoEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &mongoConfig{raw: raw}, nil
}

func (cfg *mongoConfig) Address() string {
	return fmt.Sprintf(
		"mongodb://%s:%s@%s:%s/?authSource=%s",
		cfg.raw.Username,
		cfg.raw.Password,
		cfg.raw.Host,
		cfg.raw.Port,
		cfg.raw.AuthSource,
	)
}

func (cfg *mongoConfig) DatabaseName() string {
	return cfg.raw.MongoDBName
}
