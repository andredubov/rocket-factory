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

	"github.com/andredubov/rocket-factory/order/internal/api/v1/order"
	"github.com/andredubov/rocket-factory/order/internal/repository/order/memory"
	order_v1 "github.com/andredubov/rocket-factory/shared/pkg/openapi/order/v1"
)

const (
	httpPort             = "8080"
	readHeaderTimeout    = 5 * time.Second
	shutdownTimeout      = 30 * time.Second
	paymentServerAddress = "localhost:50050"
)

func main() {
	ordersRepository := memory.NewOrderRepository()
	ordersHandler := order.NewOrderHandler(ordersRepository)

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
