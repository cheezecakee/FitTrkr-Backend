package playlist

import (
	"context"
	"database/sql"
	"log"

	"github.com/google/uuid"

	"github.com/cheezecakee/fitrkr/internal/utils/transaction"
)

type PlaylistRepo interface {
	// Playlist CRUD
	Create(ctx context.Context, playlist Playlist) (Playlist, error)
	GetByID(ctx context.Context, id int) (Playlist, error)
	GetUserPlaylists(ctx context.Context, userID uuid.UUID) ([]Playlist, error)
	Update(ctx context.Context, playlist Playlist) (Playlist, error)
	Delete(ctx context.Context, id int) error

	// Playlist with details
	GetPlaylistWithBlocks(ctx context.Context, id int) (Playlist, error)

	// Tag operations
	GetAllTags(ctx context.Context, userID uuid.UUID) ([]Tag, error)
	CreateTag(ctx context.Context, name string) (Tag, error)
	AddTagsToPlaylist(ctx context.Context, playlistID int, tagIDs []int) error
	RemoveTagsFromPlaylist(ctx context.Context, playlistID int, tagIDs []int) error
	GetPlaylistTags(ctx context.Context, playlistID int) ([]Tag, error)
}

type playlistRepo struct {
	tx transaction.BaseRepository
}

func NewPlaylistRepo(db *sql.DB) PlaylistRepo {
	return &playlistRepo{
		tx: transaction.NewBaseRepository(db),
	}
}

// Playlist CRUD Operations

const createPlaylist = `
	INSERT INTO playlists (user_id, title, description, visibility)
	VALUES ($1, $2, $3, $4)
	RETURNING id, user_id, title, description, is_active, last_worked, visibility, created_at, updated_at`

func (r *playlistRepo) Create(ctx context.Context, playlist Playlist) (Playlist, error) {
	var newPlaylist Playlist
	err := r.tx.WithTransaction(ctx, func(tx *sql.Tx) error {
		return tx.QueryRowContext(ctx, createPlaylist,
			playlist.UserID,
			playlist.Title,
			playlist.Description,
			playlist.Visibility,
		).Scan(
			&newPlaylist.ID,
			&newPlaylist.UserID,
			&newPlaylist.Title,
			&newPlaylist.Description,
			&newPlaylist.IsActive,
			&newPlaylist.LastWorked,
			&newPlaylist.Visibility,
			&newPlaylist.CreatedAt,
			&newPlaylist.UpdatedAt,
		)
	})
	if err != nil {
		log.Printf("Create playlist failed: %v", err)
		return Playlist{}, err
	}
	return newPlaylist, nil
}

const getPlaylistByID = `
	SELECT p.id, p.user_id, p.title, p.description, p.is_active, p.last_worked, 
		   p.visibility, p.created_at, p.updated_at
	FROM playlists p
	WHERE p.id = $1`

func (r *playlistRepo) GetByID(ctx context.Context, id int) (Playlist, error) {
	var playlist Playlist

	err := r.tx.DB().QueryRowContext(ctx, getPlaylistByID, id).Scan(
		&playlist.ID,
		&playlist.UserID,
		&playlist.Title,
		&playlist.Description,
		&playlist.IsActive,
		&playlist.LastWorked,
		&playlist.Visibility,
		&playlist.CreatedAt,
		&playlist.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return Playlist{}, nil
		}
		return Playlist{}, err
	}

	return playlist, nil
}

const getUserPlaylists = `
	SELECT p.id, p.user_id, p.title, p.description, p.is_active, p.last_worked, 
		   p.visibility, p.created_at, p.updated_at
	FROM playlists p
	WHERE p.user_id = $1
	ORDER BY p.updated_at DESC`

func (r *playlistRepo) GetUserPlaylists(ctx context.Context, userID uuid.UUID) ([]Playlist, error) {
	rows, err := r.tx.DB().QueryContext(ctx, getUserPlaylists, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var playlists []Playlist
	for rows.Next() {
		var playlist Playlist

		err := rows.Scan(
			&playlist.ID,
			&playlist.UserID,
			&playlist.Title,
			&playlist.Description,
			&playlist.IsActive,
			&playlist.LastWorked,
			&playlist.Visibility,
			&playlist.CreatedAt,
			&playlist.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		playlists = append(playlists, playlist)
	}

	return playlists, rows.Err()
}

const updatePlaylist = `
	UPDATE playlists 
	SET title = COALESCE(NULLIF($2, ''), title),
		description = COALESCE($3, description),
		visibility = COALESCE(NULLIF($4, ''), visibility),
		updated_at = NOW()
	WHERE id = $1 AND user_id = $5
	RETURNING id, user_id, title, description, is_active, last_worked, visibility, created_at, updated_at`

func (r *playlistRepo) Update(ctx context.Context, playlist Playlist) (Playlist, error) {
	var updatedPlaylist Playlist
	err := r.tx.WithTransaction(ctx, func(tx *sql.Tx) error {
		return tx.QueryRowContext(ctx, updatePlaylist,
			playlist.ID,
			playlist.Title,
			playlist.Description,
			playlist.Visibility,
			playlist.UserID, // For security - user can only update their own playlists
		).Scan(
			&updatedPlaylist.ID,
			&updatedPlaylist.UserID,
			&updatedPlaylist.Title,
			&updatedPlaylist.Description,
			&updatedPlaylist.IsActive,
			&updatedPlaylist.LastWorked,
			&updatedPlaylist.Visibility,
			&updatedPlaylist.CreatedAt,
			&updatedPlaylist.UpdatedAt,
		)
	})
	if err != nil {
		log.Printf("Update playlist failed for ID %d: %v", playlist.ID, err)
		return Playlist{}, err
	}
	return updatedPlaylist, nil
}

const deletePlaylist = `DELETE FROM playlists WHERE id = $1`

func (r *playlistRepo) Delete(ctx context.Context, id int) error {
	return r.tx.WithTransaction(ctx, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, deletePlaylist, id)
		return err
	})
}

// GetPlaylistWithBlocks - placeholder for now, will implement after block repo
func (r *playlistRepo) GetPlaylistWithBlocks(ctx context.Context, id int) (Playlist, error) {
	// This will be implemented once we have block repo
	return r.GetByID(ctx, id)
}
