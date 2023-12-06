package datagen

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strings"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"

	config "app/internal/config/datagen"
)

func GenerateData(ctx context.Context, cfg config.Config) error {
	if len(cfg.DSN) == 0 {
		return fmt.Errorf("passed DSN is empty")
	}
	if cfg.EmployeesCount <= 0 || cfg.EmployeesCount > 100_000_000 {
		return fmt.Errorf(
			"expected employees count to be 0 <= count <= 1_000_000, got: %d",
			cfg.EmployeesCount,
		)
	}

	db, err := sql.Open("pgx", cfg.DSN)
	if err != nil {
		return fmt.Errorf("failed to open a connection to the DB: %w", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("failed to properly close the DB connection: %v", err)
		}
	}()

	if err := createSchema(ctx, db); err != nil {
		return fmt.Errorf("failed to create the DB schema: %w", err)
	}
	if err := generateData(ctx, db, cfg.EmployeesCount); err != nil {
		return fmt.Errorf("failed to generate DB data: %w", err)
	}

	return nil
}

func createSchema(ctx context.Context, db *sql.DB) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start a transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			if !errors.Is(err, sql.ErrTxDone) {
				log.Printf("failed to rollback the transaction: %v", err)
			}
		}
	}()

	createSchemaStmts := []string{
		`CREATE TABLE IF NOT EXISTS positions(
			id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
			title VARCHAR(200) UNIQUE NOT NULL
		)`,

		`CREATE TABLE IF NOT EXISTS employees(
			id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
			first_name VARCHAR(200) NOT NULL,
			last_name VARCHAR(200) NOT NULL,
			salary NUMERIC NOT NULL,
			position INT NOT NULL REFERENCES positions(id),
			email VARCHAR(200) NOT NULL,
			CONSTRAINT employees_salary_positive_check CHECK (salary::numeric > 0)
		)`,
	}

	for _, stmt := range createSchemaStmts {
		if _, err := tx.ExecContext(ctx, stmt); err != nil {
			return fmt.Errorf("failed to execute statement `%s`: %w", stmt, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit the transaction: %w", err)
	}
	return nil
}

func generateData(ctx context.Context, db *sql.DB, empsCount int) error {
	positionsIDs, err := generatePositionsData(ctx, db)
	if err != nil {
		var pgErr *pgconn.PgError
		if !(errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation) {
			return fmt.Errorf("failed to generate the positions data: %w", err)
		}
		positionsIDs, err = fetchPositionsIDs(ctx, db)
		if err != nil {
			return fmt.Errorf("failed to fetch positions IDs: %w", err)
		}
	}
	log.Println(positionsIDs)
	if err := generateEmployeesData(ctx, db, generateEmployeesDataOpts{
		Count:        empsCount,
		PositionsIDs: positionsIDs,
	}); err != nil {
		return fmt.Errorf("failed to generate the employees data: %w", err)
	}
	return nil
}

type positionsDataIDs []int

func generatePositionsData(ctx context.Context, db *sql.DB) (positionsDataIDs, error) {
	const stmt = `INSERT INTO positions(title)
	VALUES
		('QA'),
		('Dev'),
		('PM'),
		('Architect')
	RETURNING id`
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to start a transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			if !errors.Is(err, sql.ErrTxDone) {
				log.Printf("failed to rollback the transaction: %v", err)
			}
		}
	}()

	rows, err := tx.QueryContext(ctx, stmt)
	if err != nil {
		return nil, fmt.Errorf("failed to execute the statement `%s`: %w", stmt, err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("failed to close the rows object: %v", err)
		}
	}()

	ids, err := extractPositionsIDs(rows)
	if err != nil {
		return nil, fmt.Errorf("failed to extract positions IDs: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit the transaction: %w", err)
	}
	return ids, nil
}

func fetchPositionsIDs(ctx context.Context, db *sql.DB) (positionsDataIDs, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to start a transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			if !errors.Is(err, sql.ErrTxDone) {
				log.Printf("failed to rollback the transaction: %v", err)
			}
		}
	}()

	const query = `SELECT id FROM positions`
	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to run query `%s`: %w", query, err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("failed to close the rows: %v", err)
		}
	}()

	ids, err := extractPositionsIDs(rows)
	if err != nil {
		return nil, fmt.Errorf("failed to extract positions IDs: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit the transaction: %w", err)
	}
	return ids, nil
}

func extractPositionsIDs(rows *sql.Rows) (positionsDataIDs, error) {
	var ids positionsDataIDs
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan the row: %w", err)
		}
		ids = append(ids, id)
	}
	return ids, nil
}

type generateEmployeesDataOpts struct {
	Count        int
	PositionsIDs []int
}

func generateEmployeesData(
	ctx context.Context,
	db *sql.DB,
	opts generateEmployeesDataOpts,
) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start a transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			if !errors.Is(err, sql.ErrTxDone) {
				log.Printf("failed to rollback the transaction: %v", err)
			}
		}
	}()

	const batchSize = 10000

	for opts.Count-batchSize >= 0 {
		stmt := generateEmployeesStatement(batchSize)
		fields := generateEmployeesFields(batchSize, opts.PositionsIDs)
		if _, err := tx.Exec(stmt, fields...); err != nil {
			return fmt.Errorf("failed to insert a batch: %w", err)
		}
		opts.Count -= batchSize
		log.Printf("count: %d, bs: %d", opts.Count, batchSize)
	}
	if opts.Count > 0 {
		stmt := generateEmployeesStatement(opts.Count)
		fields := generateEmployeesFields(opts.Count, opts.PositionsIDs)
		if _, err := tx.Exec(stmt, fields...); err != nil {
			return fmt.Errorf("failed to insert a batch with: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit the transaction: %w", err)
	}
	return nil
}

func generateEmployeesStatement(count int) string {
	const stmtTmpl = `INSERT INTO employees(first_name, last_name, salary, position, email)
VALUES %s`

	valuesParts := make([]string, 0, count)
	for i := 0; i < count; i++ {
		valuesParts = append(
			valuesParts,
			fmt.Sprintf(
				"($%d, $%d, $%d, $%d, $%d)",
				i*5+1, i*5+2, i*5+3, i*5+4, i*5+5,
			),
		)
	}

	return fmt.Sprintf(stmtTmpl, strings.Join(valuesParts, ","))
}

func generateEmployeesFields(count int, positionsIDs positionsDataIDs) []any {
	values := make([]any, 0, 5*count)
	for i := 0; i < count; i++ {
		values = append(
			values,
			gofakeit.FirstName(),
			gofakeit.LastName(),
			1+rand.Intn(99_999),
			positionsIDs[rand.Intn(len(positionsIDs))],
			gofakeit.Email(),
		)
	}
	return values
}
