package app

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	api "github.com/andredubov/rocket-factory/order/internal/api/v1/order"
	inventoryClient "github.com/andredubov/rocket-factory/order/internal/client/config/env/inventory"
	paymentClient "github.com/andredubov/rocket-factory/order/internal/client/config/env/payment"
	"github.com/andredubov/rocket-factory/order/internal/client/grpc/inventory/v1"
	"github.com/andredubov/rocket-factory/order/internal/client/grpc/payment/v1"
	"github.com/andredubov/rocket-factory/order/internal/config"
	"github.com/andredubov/rocket-factory/order/internal/config/env"
	"github.com/andredubov/rocket-factory/order/internal/migrator"
	"github.com/andredubov/rocket-factory/order/internal/repository/order/postgres"
	"github.com/andredubov/rocket-factory/order/internal/service"
	orders "github.com/andredubov/rocket-factory/order/internal/service/order"
	inventory_v1 "github.com/andredubov/rocket-factory/shared/pkg/proto/inventory/v1"
	payment_v1 "github.com/andredubov/rocket-factory/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	httpConfig            config.HTTPConfig
	inventoryClientConfig config.GRPCConfig
	paymentClientConfig   config.GRPCConfig
	postgresDBConfig      config.PostgresDBConfig
	inventoryClient       service.InventoryClient
	paymentClient         service.PaymentClient
	postgresDB            *pgxpool.Pool
	ordersRepository      service.OrdersRepository
	ordersService         api.OrdersService
	serverImplementation  *api.OrderImplementation
}

func NewDIContainer() *diContainer {
	return &diContainer{}
}

func (s *diContainer) HTTPConfig() config.HTTPConfig {
	if s.httpConfig == nil {
		cfg, err := env.NewHTTPConfig()
		if err != nil {
			log.Printf("failed to get http server config: %s\n", err.Error())
			return nil
		}

		s.httpConfig = cfg
	}

	return s.httpConfig
}

func (s *diContainer) InventoryGRPCConfig() config.GRPCConfig {
	if s.inventoryClientConfig == nil {
		cfg, err := inventoryClient.NewGRPCConfig()
		if err != nil {
			log.Printf("failed to create inventory grpc client config: %v\n", err)
			return nil
		}

		s.inventoryClientConfig = cfg
	}

	return s.inventoryClientConfig
}

func (s *diContainer) PaymentGRPCConfig() config.GRPCConfig {
	if s.paymentClientConfig == nil {
		cfg, err := paymentClient.NewGRPCConfig()
		if err != nil {
			log.Printf("failed to create payment grpc client config: %v\n", err)
			return nil
		}

		s.paymentClientConfig = cfg
	}

	return s.paymentClientConfig
}

func (s *diContainer) InventoryClient() service.InventoryClient {
	if s.inventoryClient == nil {
		dialOptions := []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		}

		conn, err := grpc.NewClient(s.InventoryGRPCConfig().Address(), dialOptions...)
		if err != nil {
			log.Printf("Ошибка создания клиента сервиса Inventory: %v\n", err)
			return nil
		}

		client := inventory_v1.NewInventoryServiceClient(conn)
		s.inventoryClient = inventory.NewClient(client)
	}

	return s.inventoryClient
}

func (s *diContainer) PaymentClient() service.PaymentClient {
	if s.paymentClient == nil {
		dialOptions := []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		}

		conn, err := grpc.NewClient(s.PaymentGRPCConfig().Address(), dialOptions...)
		if err != nil {
			log.Printf("Ошибка создания клиента сервиса Payment: %v\n", err)
			return nil
		}

		client := payment_v1.NewPaymentServiceClient(conn)
		s.paymentClient = payment.NewClient(client)
	}

	return s.paymentClient
}

func (s *diContainer) PostgresDBConfig() config.PostgresDBConfig {
	if s.postgresDBConfig == nil {
		cfg, err := env.NewPostgresDBConfig()
		if err != nil {
			log.Printf("failed to get MongoDB config: %s", err.Error())
			return nil
		}

		s.postgresDBConfig = cfg
	}

	return s.postgresDBConfig
}

func (s *diContainer) PostgresDatabase(ctx context.Context) *pgxpool.Pool {
	if s.postgresDB == nil {
		dbPool, err := pgxpool.New(ctx, s.PostgresDBConfig().DSN())
		if err != nil {
			log.Printf("failed to connect to database: %v\n", err)
			return nil
		}

		err = dbPool.Ping(ctx)
		if err != nil {
			log.Printf("postgres unawailable: %v\n", err)
			return nil
		}

		s.postgresDB = dbPool

		migratorRunner := migrator.NewMigrator(stdlib.OpenDBFromPool(dbPool), s.PostgresDBConfig().MigrationDirectory())
		err = migratorRunner.Up()
		if err != nil {
			log.Printf("failed to up database migration: %v\n", err)
			return nil
		}
	}

	return s.postgresDB
}

func (s *diContainer) OrdersRepository(ctx context.Context) service.OrdersRepository {
	if s.ordersRepository == nil {
		s.ordersRepository = postgres.NewOrderRepository(
			s.PostgresDatabase(ctx),
		)
	}

	return s.ordersRepository
}

func (s *diContainer) OrdersService(ctx context.Context) api.OrdersService {
	if s.ordersService == nil {
		s.ordersService = orders.NewService(
			s.OrdersRepository(ctx),
			s.PaymentClient(),
			s.InventoryClient(),
		)
	}

	return s.ordersService
}

func (s *diContainer) ServerImplementation(ctx context.Context) *api.OrderImplementation {
	if s.serverImplementation == nil {
		s.serverImplementation = api.NewOrderHandler(
			s.OrdersService(ctx),
			s.PaymentClient(),
			s.InventoryClient(),
		)
	}

	return s.serverImplementation
}
