package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	handler "github.com/andredubov/rocket-factory/order/internal/api/v1/order"
	"github.com/andredubov/rocket-factory/order/internal/client/grpc/inventory/v1"
	"github.com/andredubov/rocket-factory/order/internal/client/grpc/payment/v1"
	"github.com/andredubov/rocket-factory/order/internal/migrator"
	"github.com/andredubov/rocket-factory/order/internal/repository/order/postgres"
	"github.com/andredubov/rocket-factory/order/internal/service"
	orders "github.com/andredubov/rocket-factory/order/internal/service/order"
	"github.com/andredubov/rocket-factory/shared/pkg/config/env"
	order_v1 "github.com/andredubov/rocket-factory/shared/pkg/openapi/order/v1"
	inventory_v1 "github.com/andredubov/rocket-factory/shared/pkg/proto/inventory/v1"
	payment_v1 "github.com/andredubov/rocket-factory/shared/pkg/proto/payment/v1"
)

const (
	pingTimeout             = 5 * time.Second
	shutdownTimeout         = 30 * time.Second
	inventoryServiceAddress = "inventory-service:50051"
	paymentServiceAddress   = "payment-service:50052"
)

func main() {
	inventoryServiceClient := newInventoryServiceClient(inventoryServiceAddress)
	paymentServiceClient := newPaymentServiceClient(paymentServiceAddress)

	httpConfig, err := env.NewHTTPConfig()
	if err != nil {
		log.Printf("failed to create http server config: %v\n", err)
	}

	dbConfig, err := env.NewPostgresConfig()
	if err != nil {
		log.Printf("failed to create postgres config: %v\n", err)
	}

	ctx := context.Background()
	dbPool, err := pgxpool.New(ctx, dbConfig.DSN())
	if err != nil {
		log.Printf("failed to connect to database: %v\n", err)
		return
	}
	defer dbPool.Close()

	ctx, cancel := context.WithTimeout(ctx, pingTimeout)
	defer cancel()

	err = dbPool.Ping(ctx)
	if err != nil {
		log.Printf("postgres unawailable: %v\n", err)
		return
	}

	migratorRunner := migrator.NewMigrator(stdlib.OpenDBFromPool(dbPool), dbConfig.MigrationDirectory())
	err = migratorRunner.Up()
	if err != nil {
		log.Printf("failed to up database migration: %v\n", err)
		return
	}

	ordersRepository := postgres.NewOrderRepository(dbPool)
	ordersService := orders.NewService(ordersRepository, paymentServiceClient, inventoryServiceClient)
	ordersHandler := handler.NewOrderHandler(ordersService, paymentServiceClient, inventoryServiceClient)

	orderServer, err := order_v1.NewServer(ordersHandler)
	if err != nil {
		log.Printf("failed to create order server: %v\n", err)
		return
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Mount("/", orderServer)

	server := &http.Server{
		Addr:              httpConfig.Address(),
		Handler:           router,
		ReadHeaderTimeout: httpConfig.ReadHeaderTimeout(),
	}

	go func() {
		log.Printf("http server started on %s\n", httpConfig.Address())
		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("failed to start http server: %v\n", err)
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	<-quit

	ctx, cancel = context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("failed to shutdown http server: %v\n", err)
		return
	}

	log.Printf("http server stopped on %s\n", httpConfig.Address())
}

func newPaymentServiceClient(serviceAddress string) service.PaymentClient {
	dialOptions := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.NewClient(serviceAddress, dialOptions...)
	if err != nil {
		log.Fatalf("Ошибка создания клиента сервиса Payment: %v\n", err)
	}

	client := payment_v1.NewPaymentServiceClient(conn)
	if err != nil {
		log.Fatalf("Ошибка создания клиента сервиса Payment: %v\n", err)
	}

	return payment.NewClient(client)
}

func newInventoryServiceClient(serviceAddress string) service.InventoryClient {
	dialOptions := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.NewClient(serviceAddress, dialOptions...)
	if err != nil {
		log.Fatalf("Ошибка создания клиента сервиса Inventory: %v\n", err)
	}

	client := inventory_v1.NewInventoryServiceClient(conn)
	if err != nil {
		log.Fatalf("Ошибка создания клиента сервиса Inventory: %v\n", err)
	}

	return inventory.NewClient(client)
}
