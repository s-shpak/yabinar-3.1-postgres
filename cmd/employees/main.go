package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	config "app/internal/config/employees"
	"app/internal/core/application"
	"app/internal/infra/api"
	"app/internal/infra/store"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg := config.GetConfig()

	s, err := store.NewStore(context.Background(), store.Config{
		DSN: cfg.DSN,
	})
	if err != nil {
		return fmt.Errorf("failed to initialize a store: %w", err)
	}
	defer func() {
		if err := s.Close(); err != nil {
			log.Printf("failed to close the store: %v", err)
		}
	}()

	app := application.NewApplication(s)
	srv := api.InitServer(api.Config{Host: cfg.Host}, app)

	log.Println("Starting to accept requests")
	if err := srv.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("listen and serve has exited with an error: %w", err)
		}
	}

	return nil
}
