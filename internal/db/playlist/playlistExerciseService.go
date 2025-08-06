package playlist

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

// AddExerciseToPlaylist adds an exercise to a playlist
func (s *playlistService) AddExerciseToPlaylist(ctx context.Context, playlistID int, userID uuid.UUID, req AddExerciseToPlaylistRequest) (PlaylistExercise, error) {
	// Validate access
	if err := s.ValidatePlaylistAccess(ctx, playlistID, userID); err != nil {
		return PlaylistExercise{}, err
	}

	var blockID int

	// Handle block - create new or use existing
	if req.BlockID != nil {
		blockID = *req.BlockID
		// Verify block belongs to this playlist
		block, err := s.blockRepo.GetByID(ctx, blockID)
		if err != nil || block.PlaylistID != playlistID {
			return PlaylistExercise{}, ErrBlockNotFound
		}
	} else {
		// Create new block
		blockName := "New Block"
		if req.BlockName != nil {
			blockName = *req.BlockName
		}

		// Get next block order
		blocks, err := s.blockRepo.GetPlaylistBlocks(ctx, playlistID)
		if err != nil {
			return PlaylistExercise{}, err
		}

		newBlock := Block{
			PlaylistID:            playlistID,
			Name:                  blockName,
			BlockType:             BlockTypePlaylist,
			BlockOrder:            len(blocks) + 1,
			RestAfterBlockSeconds: 60,
		}

		createdBlock, err := s.blockRepo.Create(ctx, newBlock)
		if err != nil {
			return PlaylistExercise{}, err
		}
		blockID = createdBlock.ID
	}

	// Create exercise config
	config, err := s.configRepo.Create(ctx, req.Config)
	if err != nil {
		return PlaylistExercise{}, fmt.Errorf("failed to create config: %w", err)
	}

	// Get next exercise order within block
	blockExercises, err := s.playlistExerciseRepo.GetBlockExercises(ctx, blockID)
	if err != nil {
		return PlaylistExercise{}, err
	}

	// Create playlist exercise
	playlistExercise := PlaylistExercise{
		PlaylistID:    playlistID,
		ExerciseID:    req.ExerciseID,
		BlockID:       blockID,
		ConfigID:      config.ID,
		ExerciseOrder: len(blockExercises) + 1,
	}

	return s.playlistExerciseRepo.Create(ctx, playlistExercise)
}

// RemoveExerciseFromPlaylist removes an exercise from playlist
func (s *playlistService) RemoveExerciseFromPlaylist(ctx context.Context, exerciseID int, userID uuid.UUID) error {
	// Get exercise to validate access
	exercise, err := s.playlistExerciseRepo.GetByID(ctx, exerciseID)
	if err != nil {
		return err
	}

	if exercise.ID == 0 {
		return errors.New("exercise not found")
	}

	// Validate playlist access
	if err := s.ValidatePlaylistAccess(ctx, exercise.PlaylistID, userID); err != nil {
		return err
	}

	return s.playlistExerciseRepo.Delete(ctx, exerciseID)
}

// UpdateConfig updates exercise configuration
func (s *playlistService) UpdateConfig(ctx context.Context, exerciseID int, userID uuid.UUID, config Config) (Config, error) {
	// Get exercise to validate access
	exercise, err := s.playlistExerciseRepo.GetByID(ctx, exerciseID)
	if err != nil {
		return Config{}, err
	}

	if exercise.ID == 0 {
		return Config{}, errors.New("exercise not found")
	}

	// Validate playlist access
	if err := s.ValidatePlaylistAccess(ctx, exercise.PlaylistID, userID); err != nil {
		return Config{}, err
	}

	// Update config
	config.ID = exercise.ConfigID
	return s.configRepo.Update(ctx, config)
}

// GetPlaylistForSession returns complete playlist data for starting a workout
func (s *playlistService) GetPlaylistForSession(ctx context.Context, id int, userID uuid.UUID) (Playlist, error) {
	// Get basic playlist
	playlist, err := s.GetPlaylistByID(ctx, id, userID)
	if err != nil {
		return Playlist{}, err
	}

	// Get all blocks
	blocks, err := s.blockRepo.GetPlaylistBlocks(ctx, id)
	if err != nil {
		return Playlist{}, fmt.Errorf("failed to get blocks: %w", err)
	}

	// For each block, get exercises with configs
	for i, block := range blocks {
		exercises, err := s.playlistExerciseRepo.GetBlockExercises(ctx, block.ID)
		if err != nil {
			return Playlist{}, fmt.Errorf("failed to get exercises for block %d: %w", block.ID, err)
		}

		// Get configs for all exercises in this block
		var configIDs []int
		for _, exercise := range exercises {
			configIDs = append(configIDs, exercise.ConfigID)
		}

		configs, err := s.configRepo.GetByIDs(ctx, configIDs)
		if err != nil {
			return Playlist{}, fmt.Errorf("failed to get configs: %w", err)
		}

		// Create config map for quick lookup
		configMap := make(map[int]Config)
		for _, config := range configs {
			configMap[config.ID] = config
		}

		// Attach configs to exercises
		for j, exercise := range exercises {
			if config, exists := configMap[exercise.ConfigID]; exists {
				exercises[j].Config = &config
			}
		}

		blocks[i].Exercises = exercises
	}

	playlist.Blocks = blocks
	return playlist, nil
}
