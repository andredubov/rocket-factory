package config

type LoggerConfig interface {
	Level() string
	AsJson() bool
	EnableOTLP() bool
	OTLPEndpoint() string
	ServiceName() string
	ServiceEnvironment() string
}

type GRPCConfig interface {
	Address() string
}

type MongoDBConfig interface {
	Address() string
	DatabaseName() string
}
