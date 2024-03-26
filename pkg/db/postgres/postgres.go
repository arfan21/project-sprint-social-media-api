package dbpostgres

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/arfan21/project-sprint-social-media-api/config"
	"github.com/arfan21/project-sprint-social-media-api/pkg/logger"
	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	maxOpenConnection = 10
	connMaxLifetime   = 180
	maxIdleConns      = 30
	connMaxIdleTime   = 20
)

func NewPgxPool() (db *pgxpool.Pool, err error) {
	url := config.Get().Database.GetURL()
	// url += "?sslmode=disable"
	// if config.Get().Env != "dev" {
	// 	url += "?sslmode=verify-full&sslrootcert=ap-southeast-1-bundle.pem"
	// }
	ctx := context.Background()
	pgConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		err = fmt.Errorf("failed to parse database config: %w", err)
		return nil, err
	}

	if strings.Contains(url, "rds") {
		maxOpenConnection = 100
	}

	pgConfig.MaxConns = int32(maxOpenConnection)
	pgConfig.MaxConnIdleTime = time.Duration(connMaxIdleTime) * time.Minute
	pgConfig.HealthCheckPeriod = time.Minute
	pgConfig.ConnConfig.ConnectTimeout = 5 * time.Second
	pgConfig.MinConns = int32(maxOpenConnection) / 2
	// pgConfig.MaxConnLifetime = connMaxLifetime * time.Second

	if config.Get().Otel.EnableTracing {
		pgConfig.ConnConfig.Tracer = otelpgx.NewTracer(
			otelpgx.WithIncludeQueryParameters(),
		)
	}

	db, err = pgxpool.NewWithConfig(ctx, pgConfig)
	if err != nil {
		err = fmt.Errorf("failed to connect to database: %w", err)
		return nil, err
	}

	if err = db.Ping(ctx); err != nil {
		err = fmt.Errorf("failed to ping database: %w", err)
		return nil, err
	}

	logger.Log(ctx).Info().Msg("dbpostgres: connection established")
	return db, nil
}

func NewPgx() (db *pgx.Conn, err error) {
	dsn := config.Get().Database.GetURL()
	dsn += "?sslmode=disable"
	// if config.Get().Env == "dev" {
	// 	dsn += "?sslmode=disable"
	// } else {
	// 	dsn += "?sslmode=verify-full&sslrootcert=ap-southeast-1-bundle.pem"
	// }
	ctx := context.Background()
	pgConfig, err := pgx.ParseConfig(dsn)
	if err != nil {
		err = fmt.Errorf("failed to parse database config: %w", err)
		return nil, err
	}
	pgConfig.ConnectTimeout = 5 * time.Second

	if config.Get().Otel.EnableTracing {
		pgConfig.Tracer = otelpgx.NewTracer(
			otelpgx.WithIncludeQueryParameters(),
		)
	}

	db, err = pgx.ConnectConfig(ctx, pgConfig)
	if err != nil {
		err = fmt.Errorf("failed to connect to database: %w", err)
		return nil, err
	}

	if err = db.Ping(ctx); err != nil {
		err = fmt.Errorf("failed to ping database: %w", err)
		return nil, err
	}

	logger.Log(ctx).Info().Msg("dbpostgres: connection established")
	return db, nil

}

type Queryer interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
	CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
}
