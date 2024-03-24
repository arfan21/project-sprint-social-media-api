package dbpostgres

import (
	"context"
	"fmt"
	"time"

	"github.com/arfan21/project-sprint-social-media-api/config"
	"github.com/arfan21/project-sprint-social-media-api/pkg/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	maxOpenConnection = 100
	connMaxLifetime   = 120
	maxIdleConns      = 30
	connMaxIdleTime   = 20
)

func NewPgx() (db *pgxpool.Pool, err error) {
	dsn := config.Get().Database.GetDSN()

	if config.Get().Env == "dev" {
		dsn += "?sslmode=disable"
	} else {
		dsn += "?sslmode=verify-full&sslrootcert=ap-southeast-1-bundle.pem"
	}
	fmt.Println(dsn)
	ctx := context.Background()
	pgConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		err = fmt.Errorf("failed to parse database config: %w", err)
		return nil, err
	}
	pgConfig.MaxConns = maxOpenConnection
	pgConfig.MaxConnIdleTime = connMaxIdleTime * time.Second
	pgConfig.MaxConnLifetime = connMaxLifetime * time.Second

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

type Queryer interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
	CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
}
