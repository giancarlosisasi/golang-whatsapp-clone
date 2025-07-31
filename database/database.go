package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupDatabase(dbUrl string) *pgxpool.Pool {
	dbConfig, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create db configuration: %v\n", err)
		os.Exit(1)
	}

	dbConfig.MaxConns = 20
	dbConfig.MaxConnIdleTime = 60 * time.Second

	dbpool, err := pgxpool.NewWithConfig(context.Background(), dbConfig)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}

	err = dbpool.Ping(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to ping to database: %v\n", err)
		os.Exit(1)
	}

	return dbpool
}
