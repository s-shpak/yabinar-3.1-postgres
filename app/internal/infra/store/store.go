package store

import (
	"context"
	"fmt"

	"app/internal/core/model"
	"app/internal/infra/store/internal/db"
)

type Config struct {
	DSN string
}

type Store interface {
	GetEmployeesByName(
		ctx context.Context,
		name string,
		offsetOpts model.OffsetRequest,
	) ([]model.Employee, error)

	Close() error
}

func NewStore(ctx context.Context, cfg Config) (Store, error) {
	store, err := db.NewDB(ctx, db.Config{
		DSN: cfg.DSN,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize a DB store: %w", err)
	}
	return store, nil
}
