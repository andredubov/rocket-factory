package config

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type GRPCConfig interface {
	Address() string
}

type MongoDBConfig interface {
	Address() string
	DatabaseName() string
}
