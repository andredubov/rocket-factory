package testcontainers

// MongoDB constants
const (
	// MongoDB container constants
	MongoContainerName = "mongo"
	MongoPort          = "27017"
	// MongoDB environment variables
	MongoImageNameKey = "MONGO_IMAGE_NAME"
	MongoHostKey      = "MONGO_HOST"
	MongoPortKey      = "MONGO_PORT"
	MongoDatabaseKey  = "MONGO_INITDB_DATABASE"
	MongoUsernameKey  = "MONGO_INITDB_ROOT_USERNAME"
	MongoPasswordKey  = "MONGO_INITDB_ROOT_PASSWORD" //nolint:gosec
	MongoAuthDBKey    = "MONGO_AUTH_SOURCE"
)

// PostgresDB constants
const (
	// PostgresDB container constants
	PostgresContainerName = "postgres"
	PostgresPort          = "5432"
	// PostgresDB environment variables
	PostgresImageNameKey = "PG_IMAGE_NAME"
	PostgresHostKey      = "PG_HOST"
	PostgresPortKey      = "PG_PORT"
	PostgresDatabaseKey  = "PG_DB"
	PostgresUsernameKey  = "PG_USER"
	PostgresPasswordKey  = "PG_PASSWORD"
	PostgresSSLMode      = "PG_SSL_MODE"
	PostgresMigrationDir = "MIGRATION_DIR"
)

// Redis constants
const (
	// Redis container constants
	RedisContainerName = "redis"
	RedisPort          = "6379"
	// Redis environment variables
	RedisImageNameKey   = "REDIS_IMAGE_NAME"
	RedisHostKey        = "REDIS_HOST"
	RedisPortKey        = "REDIS_PORT"
	RedisConnTimeoutKey = "REDIS_CONNECTION_TIMEOUT"
	RedisMaxIdleKey     = "REDIS_MAX_IDLE"
	RedisIdleTimeoutKey = "REDIS_IDLE_TIMEOUT"
)
