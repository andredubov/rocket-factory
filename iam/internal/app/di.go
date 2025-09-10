package app

import (
	"context"
	"log"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"

	authAPI "github.com/andredubov/rocket-factory/iam/internal/api/v1/auth"
	userAPI "github.com/andredubov/rocket-factory/iam/internal/api/v1/user"
	"github.com/andredubov/rocket-factory/iam/internal/config"
	"github.com/andredubov/rocket-factory/iam/internal/config/env"
	"github.com/andredubov/rocket-factory/iam/internal/repository/session/redis"
	"github.com/andredubov/rocket-factory/iam/internal/repository/user/postgres"
	"github.com/andredubov/rocket-factory/iam/internal/service"
	"github.com/andredubov/rocket-factory/iam/internal/service/auth"
	"github.com/andredubov/rocket-factory/iam/internal/service/hasher"
	"github.com/andredubov/rocket-factory/iam/internal/service/user"
	"github.com/andredubov/rocket-factory/platform/pkg/cache"
	rediscache "github.com/andredubov/rocket-factory/platform/pkg/cache/redis"
	"github.com/andredubov/rocket-factory/platform/pkg/logger"
	"github.com/andredubov/rocket-factory/platform/pkg/migrator"
)

// diContainer implements the dependency container pattern
type diContainer struct {
	grpcConfig       config.GRPCConfig
	postgresDBConfig config.PostgresDBConfig
	redisConfig      config.RedisConfig

	postgresConnPool *pgxpool.Pool
	redisConnPool    *redigo.Pool
	redisClient      cache.RedisClient

	passwordHasher     service.PasswordHasher
	usersRepository    service.UsersRepository
	sessionsRepository service.SessionsRepository

	authService    authAPI.AuthService
	authAPIHandler *authAPI.AuthImplementation // GRPC service implementation

	userService    userAPI.UserService
	userAPIHandler *userAPI.UserImplementation // GRPC service implementation
}

// newDIContainer creates a new service provider instance.
func NewDIContainer() *diContainer {
	return &diContainer{}
}

// GRPCConfig loads GRPC configuration from environment variables
func (s *diContainer) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := env.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %s", err.Error())
		}
		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

func (d *diContainer) RedisConfig() config.RedisConfig {
	if d.redisConfig == nil {
		cfg, err := env.NewRedisConfig()
		if err != nil {
			log.Printf("failed to get Redis cache config: %s", err.Error())
			return nil
		}
		d.redisConfig = cfg
	}

	return d.redisConfig
}

func (d *diContainer) PostgresConfig() config.PostgresDBConfig {
	if d.postgresDBConfig == nil {
		cfg, err := env.NewPostgresDBConfig()
		if err != nil {
			log.Printf("failed to get Postgres database config: %s", err.Error())
			return nil
		}
		d.postgresDBConfig = cfg
	}

	return d.postgresDBConfig
}

func (d *diContainer) PostgresDatabase(ctx context.Context) *pgxpool.Pool {
	if d.postgresConnPool == nil {
		dbPool, err := pgxpool.New(ctx, d.PostgresConfig().DSN())
		if err != nil {
			log.Printf("failed to connect to database: %v\n", err)
			return nil
		}

		err = dbPool.Ping(ctx)
		if err != nil {
			log.Printf("postgres unawailable: %v\n", err)
			return nil
		}

		d.postgresConnPool = dbPool

		migratorRunner := migrator.NewMigrator(stdlib.OpenDBFromPool(dbPool), d.PostgresConfig().MigrationDirectory())
		err = migratorRunner.Up()
		if err != nil {
			log.Printf("failed to up database migration: %v\n", err)
			return nil
		}
	}

	return d.postgresConnPool
}

func (d *diContainer) UsersRepository(ctx context.Context) service.UsersRepository {
	if d.usersRepository == nil {
		d.usersRepository = postgres.NewUsersRepository(
			d.PostgresDatabase(ctx),
		)
	}

	return d.usersRepository
}

func (d *diContainer) RedisPool() *redigo.Pool {
	if d.redisConnPool == nil {
		d.redisConnPool = &redigo.Pool{
			MaxIdle:     config.AppConfig().RedisCache.MaxIdle(),
			IdleTimeout: config.AppConfig().RedisCache.IdleTimeout(),
			DialContext: func(ctx context.Context) (redigo.Conn, error) {
				return redigo.DialContext(ctx, "tcp", config.AppConfig().RedisCache.Address())
			},
		}
	}

	return d.redisConnPool
}

func (d *diContainer) RedisClient() cache.RedisClient {
	if d.redisClient == nil {
		d.redisClient = rediscache.NewClient(
			d.RedisPool(),
			logger.Logger(),
			config.AppConfig().RedisCache.ConnectionTimeout(),
		)
	}

	return d.redisClient
}

func (d *diContainer) SessionsRepository() service.SessionsRepository {
	if d.sessionsRepository == nil {
		d.sessionsRepository = redis.NewSessionsRepository(
			d.RedisClient(),
		)
	}

	return d.sessionsRepository
}

func (d *diContainer) PasswordHasher() service.PasswordHasher {
	if d.passwordHasher == nil {
		d.passwordHasher = hasher.NewBcryptHasher()
	}

	return d.passwordHasher
}

func (d *diContainer) AuthService(ctx context.Context) authAPI.AuthService {
	if d.authService == nil {
		d.authService = auth.NewAuthService(
			d.UsersRepository(ctx),
			d.SessionsRepository(),
			config.AppConfig().RedisCache.CacheTTL(),
			d.PasswordHasher(),
		)
	}

	return d.authService
}

func (d *diContainer) UserService(ctx context.Context) userAPI.UserService {
	if d.userService == nil {
		d.userService = user.NewUsersService(
			d.UsersRepository(ctx),
			d.PasswordHasher(),
		)
	}

	return d.userService
}

func (d *diContainer) AuthServerImplementation(ctx context.Context) *authAPI.AuthImplementation {
	if d.authAPIHandler == nil {
		d.authAPIHandler = authAPI.NewAuthImplementation(
			d.AuthService(ctx),
		)
	}

	return d.authAPIHandler
}

func (d *diContainer) UsersServerImplementation(ctx context.Context) *userAPI.UserImplementation {
	if d.userAPIHandler == nil {
		d.userAPIHandler = userAPI.NewUserImplementation(
			d.UserService(ctx),
		)
	}

	return d.userAPIHandler
}
