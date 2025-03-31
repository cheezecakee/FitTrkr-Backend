package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	u "github.com/cheezecakee/go-backend-utils/pkg/util"
	"github.com/google/uuid"

	m "github.com/cheezecakee/fitrkr/internal/models"
)

type ExSetRepo interface {
	Create(ctx context.Context, set *m.ExSet) error
	CreateBatch(ctx context.Context, sets []*m.ExSet) error
	GetByID(ctx context.Context, id uint) (*m.ExSet, error)
	GetBySessionExID(ctx context.Context, sessionExID uuid.UUID) ([]*m.ExSet, error)
	Update(ctx context.Context, set *m.ExSet) error
	Delete(ctx context.Context, id uint) error
}

type DBExSetRepo struct {
	*u.BaseRepository
}

func NewExSet(db *sql.DB) ExSetRepo {
	return &DBExSetRepo{
		u.NewBaseRepository(db),
	}
}

const createExSet = `INSERT INTO exercise_sets (session_exercise_id, set_number, reps, weight, duration, distance, notes) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

func (r *DBExSetRepo) Create(ctx context.Context, set *m.ExSet) error {
	err := r.WithTransaction(ctx, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, createExSet, set.SessionExerciseID, set.Number, set.Reps, set.Weight, set.Duration, set.Distance, set.Notes)
		return err
	})
	if err != nil {
		return err
	}
	return nil
}

const createBatchExSet = `INSERT INTO exercise_sets (session_exercise_id, set_number, reps, weight, duration, distance, notes) VALUES %s RETURNING id`

func (r *DBExSetRepo) CreateBatch(ctx context.Context, sets []*m.ExSet) error {
	if len(sets) == 0 {
		return nil
	}
	return r.WithTransaction(ctx, func(tx *sql.Tx) error {
		valueStrings := make([]string, 0, len(sets))
		valueArgs := make([]any, 0, len(sets)*7)

		for i, set := range sets {
			valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d,$%d,$%d,$%d,$%d,$%d)", i*7+1, i*7+2, i*7+3, i*7+4, i*7+5, i*7+6, i*7+7))
			valueArgs = append(valueArgs, set.SessionExerciseID, set.Number, set.Reps, set.Weight, set.Duration, set.Distance, set.Notes)
		}

		query := fmt.Sprintf(createBatchExSet, strings.Join(valueStrings, ","))

		_, err := tx.ExecContext(ctx, query, valueArgs...)
		return err
	})
}

const getExSetByID = `SELECT id, session_exercise_id, set_number, reps, weight, duration, distance, notes FROM exercise_sets WHERE id = $1`

func (r *DBExSetRepo) GetByID(ctx context.Context, id uint) (*m.ExSet, error) {
	row := r.DB.QueryRowContext(ctx, getExSetByID, id)
	exSet := &m.ExSet{}
	err := row.Scan(
		&exSet.ID,
		&exSet.SessionExerciseID,
		&exSet.Number,
		&exSet.Reps,
		&exSet.Weight,
		&exSet.Duration,
		&exSet.Distance,
		&exSet.Notes,
	)

	return exSet, err
}

const getExSetBySessionExID = `SELECT id, session_exercise_id, set_number, reps, weight, duration, distance, notes FROM exercise_sets WHERE id = $1`

func (r *DBExSetRepo) GetBySessionExID(ctx context.Context, sessionExID uuid.UUID) ([]*m.ExSet, error) {
	rows, err := r.DB.QueryContext(ctx, getExSetBySessionExID, sessionExID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exSets []*m.ExSet
	for rows.Next() {
		exSet := &m.ExSet{}
		err := rows.Scan(
			&exSet.ID,
			&exSet.SessionExerciseID,
			&exSet.Number,
			&exSet.Reps,
			&exSet.Weight,
			&exSet.Duration,
			&exSet.Distance,
			&exSet.Notes,
		)
		if err != nil {
			return nil, err
		}
		exSets = append(exSets, exSet)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return exSets, err
}

const updateExSet = `UPDATE exercise_sets
SET set_number = COALESCE($2, set_number), 
    reps = COALESCE($3, reps), 
    weight = COALESCE($4, weight), 
    duration = COALESCE($5, duration), 
    distance = COALESCE($6, distance), 
    notes = COALESCE($7, notes),
    updated_at = NOW()
WHERE id = $1`

func (r *DBExSetRepo) Update(ctx context.Context, set *m.ExSet) error {
	err := r.WithTransaction(ctx, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, updateExSet, set.ID, set.Number, set.Reps, set.Weight, set.Duration, set.Distance, set.Notes)
		return err
	})
	if err != nil {
		return err
	}
	return nil
}

const deleteExSet = `DELETE FROM exercise_sets WHERE id = $1`

func (r *DBExSetRepo) Delete(ctx context.Context, id uint) error {
	err := r.WithTransaction(ctx, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, deleteExSet, id)
		return err
	})
	if err != nil {
		return err
	}
	return nil
}
