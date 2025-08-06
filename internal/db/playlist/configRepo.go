// internal/db/playlist/exerciseConfigRepo.go
package playlist

import (
	"context"
	"database/sql"
	"log"

	"github.com/lib/pq"

	"github.com/cheezecakee/fitrkr/internal/utils/transaction"
)

type ConfigRepo interface {
	Create(ctx context.Context, config Config) (Config, error)
	GetByID(ctx context.Context, id int) (Config, error)
	Update(ctx context.Context, config Config) (Config, error)
	Delete(ctx context.Context, id int) error

	// Batch operations for efficiency
	GetByIDs(ctx context.Context, ids []int) ([]Config, error)
}

type exerciseConfigRepo struct {
	tx transaction.BaseRepository
}

func NewConfigRepo(db *sql.DB) ConfigRepo {
	return &exerciseConfigRepo{
		tx: transaction.NewBaseRepository(db),
	}
}

const createConfig = `
	INSERT INTO exercise_configs 
	(sets, reps_min, reps_max, weight, rest_seconds, tempo, 
	 duration_seconds, distance, target_pace, target_heart_rate, incline, notes)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	RETURNING id, sets, reps_min, reps_max, weight, rest_seconds, tempo,
		duration_seconds, distance, target_pace, target_heart_rate, incline, 
		notes, created_at, updated_at`

func (r *exerciseConfigRepo) Create(ctx context.Context, config Config) (Config, error) {
	var newConfig Config
	err := r.tx.WithTransaction(ctx, func(tx *sql.Tx) error {
		return tx.QueryRowContext(ctx, createConfig,
			config.Sets,
			config.RepsMin,
			config.RepsMax,
			config.Weight,
			config.RestSeconds,
			pq.Array(config.Tempo),
			config.DurationSeconds,
			config.Distance,
			config.TargetPace,
			config.TargetHeartRate,
			config.Incline,
			config.Notes,
		).Scan(
			&newConfig.ID,
			&newConfig.Sets,
			&newConfig.RepsMin,
			&newConfig.RepsMax,
			&newConfig.Weight,
			&newConfig.RestSeconds,
			pq.Array(&newConfig.Tempo),
			&newConfig.DurationSeconds,
			&newConfig.Distance,
			&newConfig.TargetPace,
			&newConfig.TargetHeartRate,
			&newConfig.Incline,
			&newConfig.Notes,
			&newConfig.CreatedAt,
			&newConfig.UpdatedAt,
		)
	})
	if err != nil {
		log.Printf("Create exercise config failed: %v", err)
		return Config{}, err
	}
	return newConfig, nil
}

const getConfigByID = `
	SELECT id, sets, reps_min, reps_max, weight, rest_seconds, tempo,
		duration_seconds, distance, target_pace, target_heart_rate, incline, 
		notes, created_at, updated_at
	FROM exercise_configs 
	WHERE id = $1`

func (r *exerciseConfigRepo) GetByID(ctx context.Context, id int) (Config, error) {
	var config Config
	err := r.tx.DB().QueryRowContext(ctx, getConfigByID, id).Scan(
		&config.ID,
		&config.Sets,
		&config.RepsMin,
		&config.RepsMax,
		&config.Weight,
		&config.RestSeconds,
		pq.Array(&config.Tempo),
		&config.DurationSeconds,
		&config.Distance,
		&config.TargetPace,
		&config.TargetHeartRate,
		&config.Incline,
		&config.Notes,
		&config.CreatedAt,
		&config.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return Config{}, nil
		}
		return Config{}, err
	}
	return config, nil
}

const updateConfig = `
	UPDATE exercise_configs 
	SET sets = COALESCE($2, sets),
		reps_min = COALESCE($3, reps_min),
		reps_max = COALESCE($4, reps_max),
		weight = COALESCE($5, weight),
		rest_seconds = COALESCE($6, rest_seconds),
		tempo = COALESCE($7, tempo),
		duration_seconds = COALESCE($8, duration_seconds),
		distance = COALESCE($9, distance),
		target_pace = COALESCE($10, target_pace),
		target_heart_rate = COALESCE($11, target_heart_rate),
		incline = COALESCE($12, incline),
		notes = COALESCE($13, notes),
		updated_at = NOW()
	WHERE id = $1
	RETURNING id, sets, reps_min, reps_max, weight, rest_seconds, tempo,
		duration_seconds, distance, target_pace, target_heart_rate, incline, 
		notes, created_at, updated_at`

func (r *exerciseConfigRepo) Update(ctx context.Context, config Config) (Config, error) {
	var updatedConfig Config
	err := r.tx.WithTransaction(ctx, func(tx *sql.Tx) error {
		return tx.QueryRowContext(ctx, updateConfig,
			config.ID,
			config.Sets,
			config.RepsMin,
			config.RepsMax,
			config.Weight,
			config.RestSeconds,
			pq.Array(config.Tempo),
			config.DurationSeconds,
			config.Distance,
			config.TargetPace,
			config.TargetHeartRate,
			config.Incline,
			config.Notes,
		).Scan(
			&updatedConfig.ID,
			&updatedConfig.Sets,
			&updatedConfig.RepsMin,
			&updatedConfig.RepsMax,
			&updatedConfig.Weight,
			&updatedConfig.RestSeconds,
			pq.Array(&updatedConfig.Tempo),
			&updatedConfig.DurationSeconds,
			&updatedConfig.Distance,
			&updatedConfig.TargetPace,
			&updatedConfig.TargetHeartRate,
			&updatedConfig.Incline,
			&updatedConfig.Notes,
			&updatedConfig.CreatedAt,
			&updatedConfig.UpdatedAt,
		)
	})
	if err != nil {
		log.Printf("Update exercise config failed for ID %d: %v", config.ID, err)
		return Config{}, err
	}
	return updatedConfig, nil
}

const deleteConfig = `DELETE FROM exercise_configs WHERE id = $1`

func (r *exerciseConfigRepo) Delete(ctx context.Context, id int) error {
	return r.tx.WithTransaction(ctx, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, deleteConfig, id)
		return err
	})
}

const getConfigsByIDs = `
	SELECT id, sets, reps_min, reps_max, weight, rest_seconds, tempo,
		duration_seconds, distance, target_pace, target_heart_rate, incline, 
		notes, created_at, updated_at
	FROM exercise_configs 
	WHERE id = ANY($1)`

func (r *exerciseConfigRepo) GetByIDs(ctx context.Context, ids []int) ([]Config, error) {
	if len(ids) == 0 {
		return []Config{}, nil
	}

	rows, err := r.tx.DB().QueryContext(ctx, getConfigsByIDs, pq.Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var configs []Config
	for rows.Next() {
		var config Config
		err := rows.Scan(
			&config.ID,
			&config.Sets,
			&config.RepsMin,
			&config.RepsMax,
			&config.Weight,
			&config.RestSeconds,
			pq.Array(&config.Tempo),
			&config.DurationSeconds,
			&config.Distance,
			&config.TargetPace,
			&config.TargetHeartRate,
			&config.Incline,
			&config.Notes,
			&config.CreatedAt,
			&config.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		configs = append(configs, config)
	}

	return configs, rows.Err()
}
