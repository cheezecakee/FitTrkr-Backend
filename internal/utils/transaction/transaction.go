// Package transaction provides transaction management utilities for FitTrkr.
package transaction

import (
	"context"
	"database/sql"
)

type BaseRepository interface {
	WithTransaction(ctx context.Context, fn func(tx *sql.Tx) error) error
	DB() *sql.DB
}

type baseRepository struct {
	db *sql.DB
}

func NewBaseRepository(db *sql.DB) BaseRepository {
	return &baseRepository{db: db}
}

func (r *baseRepository) WithTransaction(ctx context.Context, fn func(tx *sql.Tx) error) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		// Add logger later
		return err
	}

	// Ensure the transaction is properly commited or rolled back
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	// Execute the function with the transaction
	if err := fn(tx); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			// Log both errors later
		}
		return err
	}

	return tx.Commit()
}

func (r *baseRepository) DB() *sql.DB {
	return r.db
}
