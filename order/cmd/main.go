package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	handler "github.com/andredubov/rocket-factory/order/internal/api/v1/order"
	grpcClient "github.com/andredubov/rocket-factory/order/internal/client/grpc"
	"github.com/andredubov/rocket-factory/order/internal/client/grpc/payment/v1"
	"github.com/andredubov/rocket-factory/order/internal/repository/order/memory"
	orders "github.com/andredubov/rocket-factory/order/internal/service/order"
	order_v1 "github.com/andredubov/rocket-factory/shared/pkg/openapi/order/v1"
	inventory_v1 "github.com/andredubov/rocket-factory/shared/pkg/proto/inventory/v1"
	payment_v1 "github.com/andredubov/rocket-factory/shared/pkg/proto/payment/v1"
)

const (
	httpPort                = "8080"
	readHeaderTimeout       = 5 * time.Second
	shutdownTimeout         = 30 * time.Second
	paymentServiceAddress   = "localhost:50051"
	inventoryServiceAddress = "localhost:50052"
)

func main() {
	paymentServiceClient := newPaymentServiceClient(paymentServiceAddress)
	inventoryServiceClient := newInventoryServiceClient(inventoryServiceAddress)
	ordersRepository := memory.NewOrderRepository()
	ordersService := orders.NewService(ordersRepository)
	ordersHandler := handler.NewOrderHandler(ordersService, paymentServiceClient, inventoryServiceClient)

	orderServer, err := order_v1.NewServer(ordersHandler)
	if err != nil {
		log.Fatalf("failed to create order server: %v", err)
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Mount("/", orderServer)

	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", httpPort),
		Handler:           router,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	go func() {
		log.Println("server started")
		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("failed to shutdown server: %v", err)
	}

	log.Println("server stopped")
}

func newPaymentServiceClient(serviceAddress string) grpcClient.PaymentClient {
	dialOptions := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.NewClient(serviceAddress, dialOptions...)
	if err != nil {
		log.Fatalf("Ошибка создания клиента сервиса Payment: %v", err)
	}

	client := payment_v1.NewPaymentServiceClient(conn)
	if err != nil {
		log.Fatalf("Ошибка создания клиента сервиса Payment: %v", err)
	}

	return payment.NewClient(client)
}

func newInventoryServiceClient(serviceAddress string) inventory_v1.InventoryServiceClient {
	dialOptions := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.NewClient(serviceAddress, dialOptions...)
	if err != nil {
		log.Fatalf("Ошибка создания клиента сервиса Inventory: %v", err)
	}

	client := inventory_v1.NewInventoryServiceClient(conn)
	if err != nil {
		log.Fatalf("Ошибка создания клиента сервиса Inventory: %v", err)
	}

	return client
}
