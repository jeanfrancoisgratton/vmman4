package db_mgt

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

const createVMStatesTableSQL = `
CREATE TABLE IF NOT EXISTS public.vm_states (
    vms_id         BIGINT PRIMARY KEY,
    vms_name       VARCHAR(20) NOT NULL,
    vms_hypervisor VARCHAR(20) NOT NULL,
    vms_last_up    TIMESTAMPTZ,
    vms_last_down  TIMESTAMPTZ,
    vms_os         VARCHAR(30)
)`

const setVMStatesOwnerSQL = `
ALTER TABLE public.vm_states OWNER TO vmman`

const grantVMStatesSQL = `
GRANT ALL ON TABLE public.vm_states TO vmman`

// InitDB creates the vm_states table and applies ownership and privileges.
// The role executing this function must have the required PostgreSQL privileges.
func InitDB(ctx context.Context, pool *pgxpool.Pool) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("beginning database initialization transaction: %w", err)
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
			return fmt.Errorf("%s: %w", stmt.name, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("committing database initialization transaction: %w", err)
	}

	return nil
}
