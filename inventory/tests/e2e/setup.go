//go:build integration

package integration

import (
	"context"
	"fmt"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"

	"github.com/andredubov/rocket-factory/platform/pkg/logger"
	"github.com/andredubov/rocket-factory/platform/pkg/testcontainers/app"
	"github.com/andredubov/rocket-factory/platform/pkg/testcontainers/mongo"
	"github.com/andredubov/rocket-factory/platform/pkg/testcontainers/network"
	"github.com/andredubov/rocket-factory/platform/pkg/testcontainers/path"
	"github.com/andredubov/rocket-factory/platform/pkg/testcontainers/postgres"
	"github.com/andredubov/rocket-factory/platform/pkg/testcontainers/redis"
)

const (
	loggerLevelValue = "debug"

	InventoryAppContainerName           = "inventory-service-test"
	InventoryAppContainerPort           = "50051"
	InventoryMongoContainerName         = "inventory-mongo-test"
	InventoryMongoContainerInternalPort = "27017"
	InventoryDockerfile                 = "deploy/docker/inventory/Dockerfile"

	IAMAppContainerName              = "iam-service-test"
	IAMAppContainerPort              = "50053"
	IAMRedisContainerName            = "iam-redis-test"
	IAMRedisContainerInternalPort    = "6379"
	IAMPostgresContainerName         = "iam-postgres-test"
	IAMPostgresContainerInternalPort = "5432"
	IAMDockerfile                    = "deploy/docker/iam/Dockerfile"
)

type TestEnvironment struct {
	Network                 *network.Network
	InventoryMongoContainer *mongo.Container
	InventoryAppContainer   *app.Container
	IAMPostgresContainer    *postgres.Container
	IAMRedisContainer       *redis.Container
	IAMAppContainer         *app.Container
}

func customUsername(minLength int) string {
	base := gofakeit.Username()
	if len(base) >= minLength {
		return base
	}

	// Добавляем цифры или символы если username слишком короткий
	needed := minLength - len(base)
	randomSuffix := gofakeit.DigitN(uint(needed))

	return base + randomSuffix
}

func generateSimplePassword(minLength int) string {
	// Только буквы и цифры
	return gofakeit.Password(true, true, true, false, false, minLength)
}

func setupTestEnvironment(ctx context.Context) *TestEnvironment {
	logger.Info(ctx, "🚀 Setting up test environment with Testcontainers...")

	// Создаем Docker network для изоляции
	testNetwork, err := network.NewNetwork(ctx, "rocket-factory-net")
	if err != nil {
		logger.Fatal(ctx, "Failed to create network", zap.Error(err))
	}

	env := &TestEnvironment{
		Network: testNetwork,
	}

	// Запускаем инфраструктурные контейнеры
	if err := setupInfrastructureContainers(ctx, env); err != nil {
		logger.Fatal(ctx, "Failed to setup infrastructure", zap.Error(err))
	}

	// Запускаем сервисы
	if err := setupServiceContainers(ctx, env); err != nil {
		logger.Fatal(ctx, "Failed to setup services", zap.Error(err))
	}

	logger.Info(ctx, "✅ Test environment setup completed successfully")
	return env
}

func setupInfrastructureContainers(ctx context.Context, env *TestEnvironment) error {
	// Запускаем контейнеры параллельно для скорости
	errCh := make(chan error, 3)
	doneCh := make(chan bool, 3)

	go func() {
		if err := setupMongoContainer(ctx, env); err != nil {
			errCh <- fmt.Errorf("failed to setup MongoDB: %w", err)
		}
		doneCh <- true
	}()

	go func() {
		if err := setupPostgresContainer(ctx, env); err != nil {
			errCh <- fmt.Errorf("failed to setup PostgreSQL: %w", err)
		}
		doneCh <- true
	}()

	go func() {
		if err := setupRedisContainer(ctx, env); err != nil {
			errCh <- fmt.Errorf("failed to setup Redis: %w", err)
		}
		doneCh <- true
	}()

	// Ждем завершения всех горутин
	for i := 0; i < 3; i++ {
		select {
		case err := <-errCh:
			return err
		case <-doneCh:
			// Контейнер успешно запущен
		}
	}

	return nil
}

func setupServiceContainers(ctx context.Context, env *TestEnvironment) error {
	// Запускаем сервисы последовательно, т.к. они зависят от инфраструктуры
	if err := setupIAMService(ctx, env); err != nil {
		return fmt.Errorf("failed to setup IAM service: %w", err)
	}

	if err := setupInventoryService(ctx, env); err != nil {
		return fmt.Errorf("failed to setup Inventory service: %w", err)
	}

	return nil
}

func setupIAMService(ctx context.Context, env *TestEnvironment) error {
	logger.Info(ctx, "Starting IAM service container...")

	// Подготавливаем переменные окружения для IAM сервиса
	iamEnv := map[string]string{
		"GRPC_HOST":                "0.0.0.0",
		"GRPC_PORT":                IAMAppContainerPort,
		"LOGGER_LEVEL":             loggerLevelValue,
		"LOGGER_AS_JSON":           "true",
		"PG_HOST":                  IAMPostgresContainerName,
		"PG_PORT":                  IAMPostgresContainerInternalPort,
		"PG_DB":                    env.IAMPostgresContainer.Config().Database,
		"PG_USER":                  env.IAMPostgresContainer.Config().Username,
		"PG_PASSWORD":              env.IAMPostgresContainer.Config().Password,
		"PG_SSL_MODE":              env.IAMPostgresContainer.Config().SSLMode,
		"REDIS_HOST":               IAMRedisContainerName,
		"REDIS_PORT":               IAMRedisContainerInternalPort,
		"REDIS_CONNECTION_TIMEOUT": env.IAMRedisContainer.Config().ConnectionTimeout.String(),
		"REDIS_IDLE_TIMEOUT":       env.IAMRedisContainer.Config().IdleTimeout.String(),
		"REDIS_MAX_IDLE":           fmt.Sprintf("%d", env.IAMRedisContainer.Config().MaxIdle),
		"SESSION_TTL":              "24h",
	}

	projectRoot := path.GetProjectRoot()
	iamWaitStrategy := wait.ForListeningPort(nat.Port(IAMAppContainerPort + "/tcp")).WithStartupTimeout(1 * time.Minute)

	iamApp, err := app.NewContainer(
		ctx,
		app.WithName(IAMAppContainerName),
		app.WithDockerfile(projectRoot, IAMDockerfile),
		app.WithPort(IAMAppContainerPort),
		app.WithNetwork(env.Network.Name()),
		app.WithEnv(iamEnv),
		app.WithStartupWait(iamWaitStrategy),
		app.WithLogger(logger.Logger()),
	)
	if err != nil {
		return fmt.Errorf("failed to create IAM app container: %w", err)
	}

	env.IAMAppContainer = iamApp

	logger.Info(ctx, "IAM service started", zap.String("address", iamApp.Address()))

	return nil
}

func setupInventoryService(ctx context.Context, env *TestEnvironment) error {
	logger.Info(ctx, "Starting Inventory service container...")

	// Подготавливаем переменные окружения для Inventory сервиса
	inventoryEnv := map[string]string{
		"GRPC_HOST":      "0.0.0.0",
		"GRPC_PORT":      InventoryAppContainerPort,
		"LOGGER_LEVEL":   loggerLevelValue,
		"LOGGER_AS_JSON": "true",
		"MONGO_HOST":     InventoryMongoContainerName,
		"MONGO_PORT":     InventoryMongoContainerInternalPort,
		"MONGO_DATABASE": env.InventoryMongoContainer.Config().Database,
		"MONGO_USERNAME": env.InventoryMongoContainer.Config().Username,
		"MONGO_PASSWORD": env.InventoryMongoContainer.Config().Password,
		"MONGO_AUTH_SRC": env.InventoryMongoContainer.Config().AuthDB,
		"IAM_GRPC_HOST":  IAMAppContainerName, // Используем Docker DNS
		"IAM_GRPC_PORT":  IAMAppContainerPort,
	}

	projectRoot := path.GetProjectRoot()
	inventoryWaitStrategy := wait.ForListeningPort(nat.Port(InventoryAppContainerPort + "/tcp")).WithStartupTimeout(1 * time.Minute)

	inventoryApp, err := app.NewContainer(
		ctx,
		app.WithName(InventoryAppContainerName),
		app.WithDockerfile(projectRoot, InventoryDockerfile),
		app.WithPort(InventoryAppContainerPort),
		app.WithNetwork(env.Network.Name()),
		app.WithEnv(inventoryEnv),
		app.WithStartupWait(inventoryWaitStrategy),
		app.WithLogger(logger.Logger()),
	)
	if err != nil {
		return fmt.Errorf("failed to create Inventory app container: %w", err)
	}

	env.InventoryAppContainer = inventoryApp

	logger.Info(ctx, "Inventory service started", zap.String("address", inventoryApp.Address()))

	return nil
}

func setupPostgresContainer(ctx context.Context, env *TestEnvironment) error {
	logger.Info(ctx, "Starting PostgreSQL container...")

	postgresContainer, err := postgres.NewContainer(
		ctx,
		postgres.WithNetworkName(env.Network.Name()),
		postgres.WithContainerName(IAMPostgresContainerName),
		postgres.WithImageName("postgres:17.0-alpine3.20"),
		postgres.WithDatabase("iam-service-database"),
		postgres.WithAuth("iam-service-username", "iam-service-password"),
		postgres.WithSSLMode("disable"),
		postgres.WithMigrationDir("migrations"),
		postgres.WithLogger(logger.Logger()),
	)
	if err != nil {
		return fmt.Errorf("failed to create PostgreSQL container: %w", err)
	}

	env.IAMPostgresContainer = postgresContainer

	// Получаем конфигурацию для установки переменных окружения
	cfg := postgresContainer.Config()

	logger.Info(ctx, "PostgreSQL container started",
		zap.String("host", cfg.Host),
		zap.String("port", cfg.Port),
		zap.String("database", cfg.Database),
	)

	return nil
}

func setupRedisContainer(ctx context.Context, env *TestEnvironment) error {
	logger.Info(ctx, "Starting Redis container...")

	redisContainer, err := redis.NewContainer(
		ctx,
		redis.WithNetworkName(env.Network.Name()),
		redis.WithContainerName(IAMRedisContainerName),
		redis.WithImageName("redis:7.2.5-alpine3.20"),
		redis.WithConnectionTimeout(10*time.Second),
		redis.WithIdleTimeout(10*time.Second),
		redis.WithMaxIdle(10),
		redis.WithLogger(logger.Logger()),
	)
	if err != nil {
		return fmt.Errorf("failed to create Redis container: %w", err)
	}

	env.IAMRedisContainer = redisContainer

	// Получаем конфигурацию для установки переменных окружения
	cfg := redisContainer.Config()

	logger.Info(ctx, "Redis container started",
		zap.String("host", cfg.Host),
		zap.String("port", cfg.Port),
	)

	return nil
}

func setupMongoContainer(ctx context.Context, env *TestEnvironment) error {
	logger.Info(ctx, "Starting MongoDB container...")

	mongoContainer, err := mongo.NewContainer(
		ctx,
		mongo.WithNetworkName(env.Network.Name()),
		mongo.WithContainerName(InventoryMongoContainerName),
		mongo.WithImageName("mongo:7.0.5"),
		mongo.WithDatabase("inventory"),
		mongo.WithAuth("inventory_admin", "inventory_secret"),
		mongo.WithAuthDB("admin"),
		mongo.WithLogger(logger.Logger()),
	)
	if err != nil {
		return fmt.Errorf("failed to create MongoDB container: %w", err)
	}

	env.InventoryMongoContainer = mongoContainer

	// Получаем конфигурацию для установки переменных окружения
	cfg := mongoContainer.Config()

	logger.Info(ctx, "MongoDB container started",
		zap.String("host", cfg.Host),
		zap.String("port", cfg.Port),
		zap.String("database", cfg.Database),
	)

	return nil
}
