// vmman4
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: src/db_mgt/init.go
// Original timestamp : 2026.10.02 15:40:22

package db_mgt

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	ce "github.com/jeanfrancoisgratton/customError/v3"
)

// InitDB creates the vm_states table and applies ownership and privileges.
// The role executing this function must have the required PostgreSQL privileges.

func DBinit() *ce.CustomError {
	var pool *pgxpool.Pool
	var cerr *ce.CustomError
	ctx := context.Background()

	if pool, cerr = NewPool(ctx); cerr != nil {
		return cerr
	}
	tx, err := pool.Begin(ctx)
	if err != nil {
		return &ce.CustomError{Title: "Error beginning database initialization transaction", Message: err.Error()}
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	statements := []struct {
		name string
		sql  string
	}{
		{"creating vm_states table", createVMStatesTableSQL},
		{"setting vm_states owner", setVMStatesOwnerSQL},
		{"granting vm_states privileges", grantVMStatesSQL},
	}

	for _, stmt := range statements {
		if _, err := tx.Exec(ctx, stmt.sql); err != nil {
			return &ce.CustomError{Title: "Error on statement " + stmt.name, Message: err.Error()}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return &ce.CustomError{Title: "Error committing database initialization transaction", Message: err.Error()}
	}

	return nil
}
