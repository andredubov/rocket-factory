package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/andredubov/rocket-factory/order/internal/config"
	"github.com/andredubov/rocket-factory/platform/pkg/closer"
	"github.com/andredubov/rocket-factory/platform/pkg/logger"
	order_v1 "github.com/andredubov/rocket-factory/shared/pkg/openapi/order/v1"
)

type App struct {
	diContainer *diContainer
	httpServer  http.Server
	listener    net.Listener
}

func New(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	return a.runHTTPServer(ctx)
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDIContainer,
		a.initLogger,
		a.initCloser,
		a.initListener,
		a.initHTTPServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initDIContainer(_ context.Context) error {
	a.diContainer = NewDIContainer()
	return nil
}

func (a *App) initLogger(_ context.Context) error {
	return logger.Init(
		config.AppConfig().Logger.Level(),
		config.AppConfig().Logger.AsJson(),
	)
}

func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}

func (a *App) initListener(_ context.Context) error {
	listener, err := net.Listen("tcp", config.AppConfig().HTTPServer.Address())
	if err != nil {
		return err
	}

	closer.AddNamed("TCP listener", func(ctx context.Context) error {
		lerr := listener.Close()
		if lerr != nil && !errors.Is(lerr, net.ErrClosed) {
			return lerr
		}

		return nil
	})

	a.listener = listener

	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {
	orderServer, err := order_v1.NewServer(a.diContainer.ServerImplementation(ctx))
	if err != nil {
		log.Printf("failed to create order server: %v\n", err)
		return err
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Mount("/", orderServer)

	a.httpServer = http.Server{
		Addr:              a.diContainer.HTTPConfig().Address(),
		Handler:           router,
		ReadHeaderTimeout: a.diContainer.HTTPConfig().ReadHeaderTimeout(),
	}

	return nil
}

func (a *App) runHTTPServer(ctx context.Context) error {
	address := config.AppConfig().HTTPServer.Address()
	logger.Info(ctx, fmt.Sprintf("🚀 gRPC OrderService server starting on %s", address))

	return a.httpServer.Serve(a.listener) // Blocking call
}
