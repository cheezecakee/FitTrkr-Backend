package playlist

import (
	"context"
	"database/sql"
	"log"

	"github.com/cheezecakee/fitrkr/internal/utils/transaction"
)

type BlockRepo interface {
	Create(ctx context.Context, block Block) (Block, error)
	GetByID(ctx context.Context, id int) (Block, error)
	GetPlaylistBlocks(ctx context.Context, playlistID int) ([]Block, error)
	Update(ctx context.Context, block Block) (Block, error)
	Delete(ctx context.Context, id int) error

	// Reorder blocks within a playlist
	UpdateBlockOrders(ctx context.Context, playlistID int, blockOrders []BlockOrder) error
}

type BlockOrder struct {
	BlockID int `json:"block_id"`
	Order   int `json:"order"`
}

type exerciseBlockRepo struct {
	tx transaction.BaseRepository
}

func NewBlockRepo(db *sql.DB) BlockRepo {
	return &exerciseBlockRepo{
		tx: transaction.NewBaseRepository(db),
	}
}

const createBlock = `
	INSERT INTO exercise_blocks (playlist_id, name, block_type, block_order, rest_after_block_seconds)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, playlist_id, name, block_type, block_order, rest_after_block_seconds`

func (r *exerciseBlockRepo) Create(ctx context.Context, block Block) (Block, error) {
	var newBlock Block
	err := r.tx.WithTransaction(ctx, func(tx *sql.Tx) error {
		return tx.QueryRowContext(ctx, createBlock,
			block.PlaylistID,
			block.Name,
			block.BlockType,
			block.BlockOrder,
			block.RestAfterBlockSeconds,
		).Scan(
			&newBlock.ID,
			&newBlock.PlaylistID,
			&newBlock.Name,
			&newBlock.BlockType,
			&newBlock.BlockOrder,
			&newBlock.RestAfterBlockSeconds,
		)
	})
	if err != nil {
		log.Printf("Create exercise block failed: %v", err)
		return Block{}, err
	}
	return newBlock, nil
}

const getBlockByID = `
	SELECT id, playlist_id, name, block_type, block_order, rest_after_block_seconds
	FROM exercise_blocks 
	WHERE id = $1`

func (r *exerciseBlockRepo) GetByID(ctx context.Context, id int) (Block, error) {
	var block Block
	err := r.tx.DB().QueryRowContext(ctx, getBlockByID, id).Scan(
		&block.ID,
		&block.PlaylistID,
		&block.Name,
		&block.BlockType,
		&block.BlockOrder,
		&block.RestAfterBlockSeconds,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return Block{}, nil
		}
		return Block{}, err
	}
	return block, nil
}

const getPlaylistBlocks = `
	SELECT id, playlist_id, name, block_type, block_order, rest_after_block_seconds
	FROM exercise_blocks 
	WHERE playlist_id = $1
	ORDER BY block_order ASC`

func (r *exerciseBlockRepo) GetPlaylistBlocks(ctx context.Context, playlistID int) ([]Block, error) {
	rows, err := r.tx.DB().QueryContext(ctx, getPlaylistBlocks, playlistID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blocks []Block
	for rows.Next() {
		var block Block
		err := rows.Scan(
			&block.ID,
			&block.PlaylistID,
			&block.Name,
			&block.BlockType,
			&block.BlockOrder,
			&block.RestAfterBlockSeconds,
		)
		if err != nil {
			return nil, err
		}
		blocks = append(blocks, block)
	}

	return blocks, rows.Err()
}

const updateBlock = `
	UPDATE exercise_blocks 
	SET name = COALESCE(NULLIF($2, ''), name),
		block_type = COALESCE(NULLIF($3, ''), block_type),
		block_order = COALESCE($4, block_order),
		rest_after_block_seconds = COALESCE($5, rest_after_block_seconds)
	WHERE id = $1
	RETURNING id, playlist_id, name, block_type, block_order, rest_after_block_seconds`

func (r *exerciseBlockRepo) Update(ctx context.Context, block Block) (Block, error) {
	var updatedBlock Block
	err := r.tx.WithTransaction(ctx, func(tx *sql.Tx) error {
		return tx.QueryRowContext(ctx, updateBlock,
			block.ID,
			block.Name,
			block.BlockType,
			block.BlockOrder,
			block.RestAfterBlockSeconds,
		).Scan(
			&updatedBlock.ID,
			&updatedBlock.PlaylistID,
			&updatedBlock.Name,
			&updatedBlock.BlockType,
			&updatedBlock.BlockOrder,
			&updatedBlock.RestAfterBlockSeconds,
		)
	})
	if err != nil {
		log.Printf("Update exercise block failed for ID %d: %v", block.ID, err)
		return Block{}, err
	}
	return updatedBlock, nil
}

const deleteBlock = `DELETE FROM exercise_blocks WHERE id = $1`

func (r *exerciseBlockRepo) Delete(ctx context.Context, id int) error {
	return r.tx.WithTransaction(ctx, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, deleteBlock, id)
		return err
	})
}

const updateBlockOrder = `UPDATE exercise_blocks SET block_order = $2 WHERE id = $1 AND playlist_id = $3`

func (r *exerciseBlockRepo) UpdateBlockOrders(ctx context.Context, playlistID int, blockOrders []BlockOrder) error {
	return r.tx.WithTransaction(ctx, func(tx *sql.Tx) error {
		for _, order := range blockOrders {
			_, err := tx.ExecContext(ctx, updateBlockOrder, order.BlockID, order.Order, playlistID)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
