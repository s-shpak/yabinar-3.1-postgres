package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

type queryTracer struct{}

func (t *queryTracer) TraceQueryStart(
	ctx context.Context,
	_ *pgx.Conn,
	data pgx.TraceQueryStartData,
) context.Context {
	log.Printf("Running query %s (%v)", data.SQL, data.Args)
	return ctx
}

func (t *queryTracer) TraceQueryEnd(_ context.Context, _ *pgx.Conn, data pgx.TraceQueryEndData) {
	log.Printf("%v", data.CommandTag)
}
