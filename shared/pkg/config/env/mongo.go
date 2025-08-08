package env

import (
	"fmt"
	"os"

	"github.com/andredubov/rocket-factory/shared/pkg/config"
)

const (
	mongoHostEnvName       = "MG_HOST"
	mongoPortEnvName       = "MG_PORT"
	mongoUsernameEnvName   = "MG_INITDB_ROOT_USERNAME"
	mongoPasswordEnvName   = "MG_INITDB_ROOT_PASSWORD"
	mongoDatabaseEnvName   = "MG_INITDB_DATABASE"
	mongoAuthSourceEnvName = "MG_AUTH_SOURCE"
)

type mongoDBConfig struct {
	username    string
	password    string
	host        string
	port        string
	authSource  string
	mongoDBName string
}

// NewMongoDBConfig returns an instance of mongoDBConfig struct
func NewMongoDBConfig() (config.MongoDBConfig, error) {
	username := os.Getenv(mongoUsernameEnvName)
	if len(username) == 0 {
		return nil, fmt.Errorf("%s", "mongo database usernane not found")
	}

	password := os.Getenv(mongoPasswordEnvName)
	if len(password) == 0 {
		return nil, fmt.Errorf("%s", "mongo database usernane not found")
	}

	host := os.Getenv(mongoHostEnvName)
	if len(host) == 0 {
		return nil, fmt.Errorf("%s", "mongo database host not found")
	}

	port := os.Getenv(mongoPortEnvName)
	if len(port) == 0 {
		return nil, fmt.Errorf("%s", "mongo database port not found")
	}

	authSource := os.Getenv(mongoAuthSourceEnvName)
	if len(authSource) == 0 {
		return nil, fmt.Errorf("%s", "mongo database auth source not found")
	}

	mongoDBName := os.Getenv(mongoDatabaseEnvName)
	if len(mongoDBName) == 0 {
		return nil, fmt.Errorf("%s", "mongo database name not found")
	}

	return &mongoDBConfig{
		username:    username,
		password:    password,
		host:        host,
		port:        port,
		authSource:  authSource,
		mongoDBName: mongoDBName,
	}, nil
}

// Address returns Mongo database address
func (cfg *mongoDBConfig) Address() string {
	return fmt.Sprintf(
		"mongodb://%s:%s@%s:%s/?authSource=%s",
		cfg.username,
		cfg.password,
		cfg.host,
		cfg.port,
		cfg.authSource,
	)
}

func (cfg *mongoDBConfig) DatabaseName() string {
	return cfg.mongoDBName
}
