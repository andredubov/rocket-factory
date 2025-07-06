package main

import (
	"context"
	"log"

	"github.com/andredubov/rocket-factory/inventory/internal/app"
)

func main() {
	ctx := context.Background()

	application, err := app.NewApp(ctx)
	if err != nil {
		log.Fatalf("failed to init app: %s", err.Error())
	}

	err = application.Run()
	if err != nil {
		log.Fatalf("failed to run app: %s", err.Error())
	}
}
