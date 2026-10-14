package db_mgt

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	ce "github.com/jeanfrancoisgratton/customError/v3"
)

// NewPool creates and verifies a PostgreSQL connection pool.
// Set DATABASE_URL, for example:
// postgres://vmman:password@localhost:5432/vmman
func NewPool(ctx context.Context) (*pgxpool.Pool, *ce.CustomError) {
	connString := ""

	if dbe, dberr := LoadEnvironmentFile(); dberr != nil {
		return nil, dberr
	} else {
		connString = fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable",
			dbe.User, dbe.Password, dbe.Host, dbe.Port, dbe.DbName)
	}
	if connString == "" {
		return nil, &ce.CustomError{Title: "DATABASE_URL is not set"}
	}

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, &ce.CustomError{Title: "Error parsing database configuration", Message: err.Error()}
	}

	// Adjust these defaults to suit the application.
	config.MinConns = 1
	config.MaxConns = 5

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, &ce.CustomError{Title: "Error creating PostgreSQL connection pool", Message: err.Error()}
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, &ce.CustomError{Title: "Error cconnecting to PostgreSQL", Message: err.Error()}
	}

	return pool, nil
}
