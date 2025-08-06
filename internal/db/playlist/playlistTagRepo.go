package playlist

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
)

const getAllTags = `SELECT id, name FROM tags ORDER BY name`

func (r *playlistRepo) GetAllTags(ctx context.Context, userID uuid.UUID) ([]Tag, error) {
	rows, err := r.tx.DB().QueryContext(ctx, getAllTags)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []Tag
	for rows.Next() {
		var tag Tag
		if err := rows.Scan(&tag.ID, &tag.Name); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, rows.Err()
}

const createTag = `INSERT INTO tags (name) VALUES ($1) RETURNING id, name`

func (r *playlistRepo) CreateTag(ctx context.Context, name string) (Tag, error) {
	var tag Tag
	err := r.tx.WithTransaction(ctx, func(tx *sql.Tx) error {
		return tx.QueryRowContext(ctx, createTag, name).Scan(&tag.ID, &tag.Name)
	})
	return tag, err
}

const addTagsToPlaylist = `INSERT INTO playlist_tags (playlist_id, tag_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`

func (r *playlistRepo) AddTagsToPlaylist(ctx context.Context, playlistID int, tagIDs []int) error {
	return r.tx.WithTransaction(ctx, func(tx *sql.Tx) error {
		for _, tagID := range tagIDs {
			_, err := tx.ExecContext(ctx, addTagsToPlaylist, playlistID, tagID)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

const removeTagsFromPlaylist = `DELETE FROM playlist_tags WHERE playlist_id = $1 AND tag_id = $2`

func (r *playlistRepo) RemoveTagsFromPlaylist(ctx context.Context, playlistID int, tagIDs []int) error {
	return r.tx.WithTransaction(ctx, func(tx *sql.Tx) error {
		for _, tagID := range tagIDs {
			_, err := tx.ExecContext(ctx, removeTagsFromPlaylist, playlistID, tagID)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

const getPlaylistTags = `
	SELECT t.id, t.name 
	FROM tags t 
	JOIN playlist_tags pt ON t.id = pt.tag_id 
	WHERE pt.playlist_id = $1`

func (r *playlistRepo) GetPlaylistTags(ctx context.Context, playlistID int) ([]Tag, error) {
	rows, err := r.tx.DB().QueryContext(ctx, getPlaylistTags, playlistID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []Tag
	for rows.Next() {
		var tag Tag
		if err := rows.Scan(&tag.ID, &tag.Name); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, rows.Err()
}
