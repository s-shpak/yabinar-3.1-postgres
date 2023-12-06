package main

import (
	"context"
	"fmt"
	"log"

	config "app/internal/config/datagen"
	"app/internal/datagen"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg := config.GetConfig()
	if err := datagen.GenerateData(context.Background(), cfg); err != nil {
		return fmt.Errorf("failed to generate the DB data: %w", err)
	}
	return nil
}
