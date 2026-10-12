package db_mgt

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// UpdateLastUp records the last time a VM started.
// Pass a time.Time value directly; pgx encodes it as TIMESTAMPTZ.
func UpdateLastUp(
	ctx context.Context,
	pool *pgxpool.Pool,
	vmID int64,
	lastUp time.Time,
) error {
	const query = `
UPDATE public.vm_states
SET vms_last_up = $1
WHERE vms_id = $2`

	tag, err := pool.Exec(ctx, query, lastUp, vmID)
	if err != nil {
		return fmt.Errorf("updating VM %d last-up time: %w", vmID, err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("VM %d not found", vmID)
	}

	return nil
}

// UpdateLastDown records the last time a VM stopped.
func UpdateLastDown(
	ctx context.Context,
	pool *pgxpool.Pool,
	vmID int64,
	lastDown time.Time,
) error {
	const query = `
UPDATE public.vm_states
SET vms_last_down = $1
WHERE vms_id = $2`

	tag, err := pool.Exec(ctx, query, lastDown, vmID)
	if err != nil {
		return fmt.Errorf("updating VM %d last-down time: %w", vmID, err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("VM %d not found", vmID)
	}

	return nil
}
