package playlist

import (
	"context"

	"github.com/google/uuid"
)

// CreateBlock creates a new exercise block
func (s *playlistService) CreateBlock(ctx context.Context, playlistID int, userID uuid.UUID, blockName string, blockType string) (Block, error) {
	// Validate access
	if err := s.ValidatePlaylistAccess(ctx, playlistID, userID); err != nil {
		return Block{}, err
	}

	// Validate block type
	if !isValidBlockType(blockType) {
		return Block{}, ErrInvalidBlockType
	}

	// Get next block order
	blocks, err := s.blockRepo.GetPlaylistBlocks(ctx, playlistID)
	if err != nil {
		return Block{}, err
	}

	newBlock := Block{
		PlaylistID:            playlistID,
		Name:                  blockName,
		BlockType:             BlockType(blockType),
		BlockOrder:            len(blocks) + 1,
		RestAfterBlockSeconds: 60,
	}

	return s.blockRepo.Create(ctx, newBlock)
}

// UpdateBlockOrder reorders blocks within a playlist
func (s *playlistService) UpdateBlockOrder(ctx context.Context, playlistID int, userID uuid.UUID, blockOrders []BlockOrder) error {
	// Validate access
	if err := s.ValidatePlaylistAccess(ctx, playlistID, userID); err != nil {
		return err
	}

	return s.blockRepo.UpdateBlockOrders(ctx, playlistID, blockOrders)
}

// GetAllTags returns all available tags
func (s *playlistService) GetAllTags(ctx context.Context, userID uuid.UUID) ([]Tag, error) {
	return s.playlistRepo.GetAllTags(ctx, userID)
}
