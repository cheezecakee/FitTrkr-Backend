package playlist

import (
	"context"
	"database/sql"
	"log"

	"github.com/cheezecakee/fitrkr/internal/utils/transaction"
)

type PlaylistExerciseRepo interface {
	Create(ctx context.Context, playlistExercise PlaylistExercise) (PlaylistExercise, error)
	GetByID(ctx context.Context, id int) (PlaylistExercise, error)
	GetBlockExercises(ctx context.Context, blockID int) ([]PlaylistExercise, error)
	GetPlaylistExercises(ctx context.Context, playlistID int) ([]PlaylistExercise, error)
	Update(ctx context.Context, playlistExercise PlaylistExercise) (PlaylistExercise, error)
	Delete(ctx context.Context, id int) error

	// Reorder exercises within a block
	UpdateExerciseOrders(ctx context.Context, blockID int, exerciseOrders []ExerciseOrder) error

	// Move exercise to different block
	MoveExerciseToBlock(ctx context.Context, exerciseID int, newBlockID int, newOrder int) error
}

type ExerciseOrder struct {
	ExerciseID int `json:"exercise_id"`
	Order      int `json:"order"`
}

type playlistExerciseRepo struct {
	tx transaction.BaseRepository
}

func NewPlaylistExerciseRepo(db *sql.DB) PlaylistExerciseRepo {
	return &playlistExerciseRepo{
		tx: transaction.NewBaseRepository(db),
	}
}

const createPlaylistExercise = `
	INSERT INTO playlist_exercises (playlist_id, exercise_id, block_id, config_id, exercise_order)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, playlist_id, exercise_id, block_id, config_id, exercise_order, created_at, updated_at`

func (r *playlistExerciseRepo) Create(ctx context.Context, playlistExercise PlaylistExercise) (PlaylistExercise, error) {
	var newPlaylistExercise PlaylistExercise
	err := r.tx.WithTransaction(ctx, func(tx *sql.Tx) error {
		return tx.QueryRowContext(ctx, createPlaylistExercise,
			playlistExercise.PlaylistID,
			playlistExercise.ExerciseID,
			playlistExercise.BlockID,
			playlistExercise.ConfigID,
			playlistExercise.ExerciseOrder,
		).Scan(
			&newPlaylistExercise.ID,
			&newPlaylistExercise.PlaylistID,
			&newPlaylistExercise.ExerciseID,
			&newPlaylistExercise.BlockID,
			&newPlaylistExercise.ConfigID,
			&newPlaylistExercise.ExerciseOrder,
			&newPlaylistExercise.CreatedAt,
			&newPlaylistExercise.UpdatedAt,
		)
	})
	if err != nil {
		log.Printf("Create playlist exercise failed: %v", err)
		return PlaylistExercise{}, err
	}
	return newPlaylistExercise, nil
}

const getPlaylistExerciseByID = `
	SELECT pe.id, pe.playlist_id, pe.exercise_id, pe.block_id, pe.config_id, 
		   pe.exercise_order, pe.created_at, pe.updated_at, e.name
	FROM playlist_exercises pe
	JOIN exercises e ON pe.exercise_id = e.id
	WHERE pe.id = $1`

func (r *playlistExerciseRepo) GetByID(ctx context.Context, id int) (PlaylistExercise, error) {
	var playlistExercise PlaylistExercise
	err := r.tx.DB().QueryRowContext(ctx, getPlaylistExerciseByID, id).Scan(
		&playlistExercise.ID,
		&playlistExercise.PlaylistID,
		&playlistExercise.ExerciseID,
		&playlistExercise.BlockID,
		&playlistExercise.ConfigID,
		&playlistExercise.ExerciseOrder,
		&playlistExercise.CreatedAt,
		&playlistExercise.UpdatedAt,
		&playlistExercise.ExerciseName,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return PlaylistExercise{}, nil
		}
		return PlaylistExercise{}, err
	}
	return playlistExercise, nil
}

const getBlockExercises = `
	SELECT pe.id, pe.playlist_id, pe.exercise_id, pe.block_id, pe.config_id, 
		   pe.exercise_order, pe.created_at, pe.updated_at, e.name
	FROM playlist_exercises pe
	JOIN exercises e ON pe.exercise_id = e.id
	WHERE pe.block_id = $1
	ORDER BY pe.exercise_order ASC`

func (r *playlistExerciseRepo) GetBlockExercises(ctx context.Context, blockID int) ([]PlaylistExercise, error) {
	rows, err := r.tx.DB().QueryContext(ctx, getBlockExercises, blockID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exercises []PlaylistExercise
	for rows.Next() {
		var exercise PlaylistExercise
		err := rows.Scan(
			&exercise.ID,
			&exercise.PlaylistID,
			&exercise.ExerciseID,
			&exercise.BlockID,
			&exercise.ConfigID,
			&exercise.ExerciseOrder,
			&exercise.CreatedAt,
			&exercise.UpdatedAt,
			&exercise.ExerciseName,
		)
		if err != nil {
			return nil, err
		}
		exercises = append(exercises, exercise)
	}

	return exercises, rows.Err()
}

const getPlaylistExercises = `
	SELECT pe.id, pe.playlist_id, pe.exercise_id, pe.block_id, pe.config_id, 
		   pe.exercise_order, pe.created_at, pe.updated_at, e.name
	FROM playlist_exercises pe
	JOIN exercises e ON pe.exercise_id = e.id
	WHERE pe.playlist_id = $1
	ORDER BY pe.block_id, pe.exercise_order ASC`

func (r *playlistExerciseRepo) GetPlaylistExercises(ctx context.Context, playlistID int) ([]PlaylistExercise, error) {
	rows, err := r.tx.DB().QueryContext(ctx, getPlaylistExercises, playlistID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exercises []PlaylistExercise
	for rows.Next() {
		var exercise PlaylistExercise
		err := rows.Scan(
			&exercise.ID,
			&exercise.PlaylistID,
			&exercise.ExerciseID,
			&exercise.BlockID,
			&exercise.ConfigID,
			&exercise.ExerciseOrder,
			&exercise.CreatedAt,
			&exercise.UpdatedAt,
			&exercise.ExerciseName,
		)
		if err != nil {
			return nil, err
		}
		exercises = append(exercises, exercise)
	}

	return exercises, rows.Err()
}

const updatePlaylistExercise = `
	UPDATE playlist_exercises 
	SET block_id = COALESCE($2, block_id),
		config_id = COALESCE($3, config_id),
		exercise_order = COALESCE($4, exercise_order),
		updated_at = NOW()
	WHERE id = $1
	RETURNING id, playlist_id, exercise_id, block_id, config_id, exercise_order, created_at, updated_at`

func (r *playlistExerciseRepo) Update(ctx context.Context, playlistExercise PlaylistExercise) (PlaylistExercise, error) {
	var updatedExercise PlaylistExercise
	err := r.tx.WithTransaction(ctx, func(tx *sql.Tx) error {
		return tx.QueryRowContext(ctx, updatePlaylistExercise,
			playlistExercise.ID,
			playlistExercise.BlockID,
			playlistExercise.ConfigID,
			playlistExercise.ExerciseOrder,
		).Scan(
			&updatedExercise.ID,
			&updatedExercise.PlaylistID,
			&updatedExercise.ExerciseID,
			&updatedExercise.BlockID,
			&updatedExercise.ConfigID,
			&updatedExercise.ExerciseOrder,
			&updatedExercise.CreatedAt,
			&updatedExercise.UpdatedAt,
		)
	})
	if err != nil {
		log.Printf("Update playlist exercise failed for ID %d: %v", playlistExercise.ID, err)
		return PlaylistExercise{}, err
	}
	return updatedExercise, nil
}

const deletePlaylistExercise = `DELETE FROM playlist_exercises WHERE id = $1`

func (r *playlistExerciseRepo) Delete(ctx context.Context, id int) error {
	return r.tx.WithTransaction(ctx, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, deletePlaylistExercise, id)
		return err
	})
}

const updateExerciseOrder = `UPDATE playlist_exercises SET exercise_order = $2 WHERE id = $1 AND block_id = $3`

func (r *playlistExerciseRepo) UpdateExerciseOrders(ctx context.Context, blockID int, exerciseOrders []ExerciseOrder) error {
	return r.tx.WithTransaction(ctx, func(tx *sql.Tx) error {
		for _, order := range exerciseOrders {
			_, err := tx.ExecContext(ctx, updateExerciseOrder, order.ExerciseID, order.Order, blockID)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

const moveExerciseToBlock = `
	UPDATE playlist_exercises 
	SET block_id = $2, exercise_order = $3, updated_at = NOW() 
	WHERE id = $1`

func (r *playlistExerciseRepo) MoveExerciseToBlock(ctx context.Context, exerciseID int, newBlockID int, newOrder int) error {
	return r.tx.WithTransaction(ctx, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, moveExerciseToBlock, exerciseID, newBlockID, newOrder)
		return err
	})
}
