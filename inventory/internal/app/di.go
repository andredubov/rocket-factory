package app

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	api "github.com/andredubov/rocket-factory/inventory/internal/api/v1/inventory"
	"github.com/andredubov/rocket-factory/inventory/internal/config"
	"github.com/andredubov/rocket-factory/inventory/internal/config/env"
	mongodb "github.com/andredubov/rocket-factory/inventory/internal/repository/part/mongo"
	"github.com/andredubov/rocket-factory/inventory/internal/service"
	"github.com/andredubov/rocket-factory/inventory/internal/service/inventory"
	"github.com/andredubov/rocket-factory/platform/pkg/closer"
	"github.com/andredubov/rocket-factory/platform/pkg/logger"
	middlewaregrpc "github.com/andredubov/rocket-factory/platform/pkg/middleware/grpc"
	auth_v1 "github.com/andredubov/rocket-factory/shared/pkg/proto/auth/v1"
)

// diContainer implements the dependency container pattern
type diContainer struct {
	inventoryRepository  service.InventoryRepository
	inventoryService     api.InventoryService
	grpcConfig           config.GRPCConfig
	mongoDBConfig        config.MongoDBConfig
	mongoDB              *mongo.Database
	serverImplementation *api.InventoryImplementation // GRPC service implementation
	authClient           auth_v1.AuthServiceClient
	authInterceptor      *middlewaregrpc.AuthInterceptor
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

// MongoDBConfig loads MongoDB configuration from environment variables
func (s *diContainer) MongoDBConfig() config.MongoDBConfig {
	if s.mongoDBConfig == nil {
		cfg, err := env.NewMongoDBConfig()
		if err != nil {
			log.Fatalf("failed to get MongoDB config: %s", err.Error())
		}
		s.mongoDBConfig = cfg
	}

	return s.mongoDBConfig
}

// MongoDatabase creates an instance of database client
func (s *diContainer) MongoDatabase(ctx context.Context) *mongo.Database {
	if s.mongoDB == nil {
		URI := s.MongoDBConfig().Address()
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(URI))
		if err != nil {
			log.Printf("failed to connect MongoDB: %s\n", err)
			return nil
		}

		err = client.Ping(ctx, nil)
		if err != nil {
			log.Printf("failed to ping MongoDB: %v\n", err)
			return nil
		}

		s.mongoDB = client.Database(s.MongoDBConfig().DatabaseName())
	}

	return s.mongoDB
}

// InventoryRepository provides access to inventory data
func (s *diContainer) InventoryRepository(ctx context.Context) service.InventoryRepository {
	if s.inventoryRepository == nil {
		s.inventoryRepository = mongodb.NewInventoryRepository(
			ctx,
			s.MongoDatabase(ctx),
		)
	}

	return s.inventoryRepository
}

// InventoryService provides access to inventory service layer
func (s *diContainer) InventoryService(ctx context.Context) api.InventoryService {
	if s.inventoryService == nil {
		s.inventoryService = inventory.NewService(
			s.InventoryRepository(ctx),
		)
	}

	return s.inventoryService
}

// ServerImplementation creates GRPC service handler
func (s *diContainer) ServerImplementation(ctx context.Context) *api.InventoryImplementation {
	if s.serverImplementation == nil {
		inventoryService := s.InventoryService(ctx)
		s.serverImplementation = api.NewInventoryImplementation(inventoryService)
	}

	return s.serverImplementation
}

func (d *diContainer) AuthClient(ctx context.Context) auth_v1.AuthServiceClient {
	if d.authClient == nil {
		iamAddress := config.AppConfig().IAMClient.Address()
		if iamAddress == "" {
			logger.Error(ctx, "IAM client address is empty")
			return nil
		}

		connectCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		iamConn, err := grpc.NewClient(
			iamAddress,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			logger.Error(connectCtx, "Failed to connect to IAM service", zap.Error(err))
			return nil
		}

		closer.AddNamed("IAM grpc connection", func(ctx context.Context) error {
			return iamConn.Close()
		})

		d.authClient = auth_v1.NewAuthServiceClient(iamConn)
	}

	return d.authClient
}

func (d *diContainer) AuthInterceptor(ctx context.Context) *middlewaregrpc.AuthInterceptor {
	if d.authInterceptor == nil {
		d.authInterceptor = middlewaregrpc.NewAuthInterceptor(
			d.AuthClient(ctx),
		)
	}

	return d.authInterceptor
}
