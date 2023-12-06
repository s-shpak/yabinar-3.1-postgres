package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"app/internal/core/model"
)

type Config struct {
	DSN string
}

type DB struct {
	pool *pgxpool.Pool
}

func NewDB(ctx context.Context, cfg Config) (*DB, error) {
	pool, err := initPool(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize a connection pool: %w", err)
	}
	return &DB{
		pool: pool,
	}, nil
}

func initPool(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {
	poolCfg, err := pgxpool.ParseConfig(cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the DSN: %w", err)
	}
	poolCfg.ConnConfig.Tracer = &queryTracer{}
	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize a connection pool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping the DB: %w", err)
	}
	return pool, nil
}

func (db *DB) GetEmployeesByName(
	ctx context.Context,
	name string,
	offsetOpts model.OffsetRequest,
) ([]model.Employee, error) {
	const queryStmt = `SELECT id, first_name, last_name, salary, position, email
FROM employees
WHERE last_name LIKE CONCAT('%%', $1::text, '%%') AND id > $2
LIMIT $3`
	rows, err := db.pool.Query(ctx, queryStmt, fmt.Sprintf("%%%s%%", name), offsetOpts.LastID, offsetOpts.Limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query DB: %w", err)
	}
	defer rows.Close()

	emps := make([]model.Employee, 0, offsetOpts.Limit)
	for rows.Next() {
		var emp model.Employee
		if err := rows.Scan(
			&emp.ID,
			&emp.FirstName,
			&emp.LastName,
			&emp.Salary,
			&emp.PositionID,
			&emp.Email,
		); err != nil {
			return nil, fmt.Errorf("failed to scan a response row: %w", err)
		}
		emps = append(emps, emp)
	}
	return emps, nil
}

func (db *DB) Close() error {
	db.pool.Close()
	return nil
}
