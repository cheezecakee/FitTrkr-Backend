package playlist

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
)

var (
	ErrPlaylistNotFound   = errors.New("playlist not found")
	ErrPlaylistExists     = errors.New("playlist with this title already exists")
	ErrUnauthorizedAccess = errors.New("unauthorized access to playlist")
	ErrBlockNotFound      = errors.New("exercise block not found")
	ErrInvalidBlockType   = errors.New("invalid block type")
	ErrConfigNotFound     = errors.New("exercise config not found")
)

type PlaylistService interface {
	// Playlist operations
	CreatePlaylist(ctx context.Context, userID uuid.UUID, req CreatePlaylistRequest) (Playlist, error)
	GetPlaylistByID(ctx context.Context, id int, userID uuid.UUID) (Playlist, error)
	GetUserPlaylists(ctx context.Context, userID uuid.UUID) ([]PlaylistWithDetails, error)
	UpdatePlaylist(ctx context.Context, id int, userID uuid.UUID, req UpdatePlaylistRequest) (Playlist, error)
	DeletePlaylist(ctx context.Context, id int, userID uuid.UUID) error

	// Full playlist with all exercises (for starting a session)
	GetPlaylistForSession(ctx context.Context, id int, userID uuid.UUID) (Playlist, error)

	// Exercise management
	AddExerciseToPlaylist(ctx context.Context, playlistID int, userID uuid.UUID, req AddExerciseToPlaylistRequest) (PlaylistExercise, error)
	RemoveExerciseFromPlaylist(ctx context.Context, exerciseID int, userID uuid.UUID) error
	UpdateConfig(ctx context.Context, exerciseID int, userID uuid.UUID, config Config) (Config, error)

	// Block management
	CreateBlock(ctx context.Context, playlistID int, userID uuid.UUID, blockName string, blockType string) (Block, error)
	UpdateBlockOrder(ctx context.Context, playlistID int, userID uuid.UUID, blockOrders []BlockOrder) error

	// Utility methods
	GetAllTags(ctx context.Context, userID uuid.UUID) ([]Tag, error)

	// Validation helpers
	ValidatePlaylistAccess(ctx context.Context, playlistID int, userID uuid.UUID) error
}

type playlistService struct {
	playlistRepo         PlaylistRepo
	blockRepo            BlockRepo
	playlistExerciseRepo PlaylistExerciseRepo
	configRepo           ConfigRepo
}

func NewPlaylistService(
	playlistRepo PlaylistRepo,
	blockRepo BlockRepo,
	playlistExerciseRepo PlaylistExerciseRepo,
	configRepo ConfigRepo,
) PlaylistService {
	return &playlistService{
		playlistRepo:         playlistRepo,
		blockRepo:            blockRepo,
		playlistExerciseRepo: playlistExerciseRepo,
		configRepo:           configRepo,
	}
}

// CreatePlaylist creates a new playlist with default block
func (s *playlistService) CreatePlaylist(ctx context.Context, userID uuid.UUID, req CreatePlaylistRequest) (Playlist, error) {
	// Set defaults
	if req.Visibility == "" {
		req.Visibility = string(VisibilityPrivate)
	}

	playlist := Playlist{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Visibility:  Visibility(req.Visibility),
	}

	// Create playlist
	createdPlaylist, err := s.playlistRepo.Create(ctx, playlist)
	if err != nil {
		if isUniqueConstraintError(err) {
			return Playlist{}, ErrPlaylistExists
		}
		return Playlist{}, fmt.Errorf("failed to create playlist: %w", err)
	}

	// Create default block
	defaultBlock := Block{
		PlaylistID:            createdPlaylist.ID,
		Name:                  "Playlist",
		BlockType:             BlockTypePlaylist,
		BlockOrder:            1,
		RestAfterBlockSeconds: 60,
	}

	_, err = s.blockRepo.Create(ctx, defaultBlock)
	if err != nil {
		log.Printf("Failed to create default block for playlist %d: %v", createdPlaylist.ID, err)
		// Continue - playlist is created, user can add blocks later
	}

	// Add tags if provided
	if len(req.TagIDs) > 0 {
		err = s.playlistRepo.AddTagsToPlaylist(ctx, createdPlaylist.ID, req.TagIDs)
		if err != nil {
			log.Printf("Failed to add tags to playlist %d: %v", createdPlaylist.ID, err)
			// Continue - playlist is created
		}
	}

	return createdPlaylist, nil
}

// GetPlaylistByID returns basic playlist info
func (s *playlistService) GetPlaylistByID(ctx context.Context, id int, userID uuid.UUID) (Playlist, error) {
	playlist, err := s.playlistRepo.GetByID(ctx, id)
	if err != nil {
		return Playlist{}, err
	}

	if playlist.ID == 0 {
		return Playlist{}, ErrPlaylistNotFound
	}

	// Check access
	if playlist.UserID != userID && playlist.Visibility == VisibilityPrivate {
		return Playlist{}, ErrUnauthorizedAccess
	}

	// Get tags
	tags, err := s.playlistRepo.GetPlaylistTags(ctx, id)
	if err != nil {
		log.Printf("Failed to get tags for playlist %d: %v", id, err)
	} else {
		playlist.Tags = tags
	}

	return playlist, nil
}

// GetUserPlaylists returns all playlists for a user with summary info
func (s *playlistService) GetUserPlaylists(ctx context.Context, userID uuid.UUID) ([]PlaylistWithDetails, error) {
	playlists, err := s.playlistRepo.GetUserPlaylists(ctx, userID)
	if err != nil {
		return nil, err
	}

	var playlistsWithDetails []PlaylistWithDetails
	for _, playlist := range playlists {
		// Get exercise count for each playlist
		exercises, err := s.playlistExerciseRepo.GetPlaylistExercises(ctx, playlist.ID)
		if err != nil {
			log.Printf("Failed to get exercise count for playlist %d: %v", playlist.ID, err)
			continue
		}

		// Get block count
		blocks, err := s.blockRepo.GetPlaylistBlocks(ctx, playlist.ID)
		if err != nil {
			log.Printf("Failed to get block count for playlist %d: %v", playlist.ID, err)
			continue
		}

		// Get tags
		tags, err := s.playlistRepo.GetPlaylistTags(ctx, playlist.ID)
		if err == nil {
			playlist.Tags = tags
		}

		playlistWithDetails := PlaylistWithDetails{
			Playlist:       playlist,
			TotalExercises: len(exercises),
			TotalBlocks:    len(blocks),
		}

		playlistsWithDetails = append(playlistsWithDetails, playlistWithDetails)
	}

	return playlistsWithDetails, nil
}

// UpdatePlaylist updates playlist details
func (s *playlistService) UpdatePlaylist(ctx context.Context, id int, userID uuid.UUID, req UpdatePlaylistRequest) (Playlist, error) {
	// Validate access
	if err := s.ValidatePlaylistAccess(ctx, id, userID); err != nil {
		return Playlist{}, err
	}

	// Build update struct with only provided fields
	updatePlaylist := Playlist{
		ID:     id,
		UserID: userID,
	}

	if req.Title != nil {
		updatePlaylist.Title = *req.Title
	}
	if req.Description != nil {
		updatePlaylist.Description = req.Description
	}
	if req.Visibility != nil {
		updatePlaylist.Visibility = Visibility(*req.Visibility)
	}

	updatedPlaylist, err := s.playlistRepo.Update(ctx, updatePlaylist)
	if err != nil {
		if isUniqueConstraintError(err) {
			return Playlist{}, ErrPlaylistExists
		}
		return Playlist{}, fmt.Errorf("failed to update playlist: %w", err)
	}

	// Update tags if provided
	if req.TagIDs != nil {
		// Remove existing tags and add new ones
		existingTags, err := s.playlistRepo.GetPlaylistTags(ctx, id)
		if err == nil && len(existingTags) > 0 {
			var existingTagIDs []int
			for _, tag := range existingTags {
				existingTagIDs = append(existingTagIDs, tag.ID)
			}
			s.playlistRepo.RemoveTagsFromPlaylist(ctx, id, existingTagIDs)
		}

		if len(req.TagIDs) > 0 {
			err = s.playlistRepo.AddTagsToPlaylist(ctx, id, req.TagIDs)
			if err != nil {
				log.Printf("Failed to update tags for playlist %d: %v", id, err)
			}
		}
	}

	return updatedPlaylist, nil
}

// DeletePlaylist removes a playlist and all associated data
func (s *playlistService) DeletePlaylist(ctx context.Context, id int, userID uuid.UUID) error {
	// Validate access
	if err := s.ValidatePlaylistAccess(ctx, id, userID); err != nil {
		return err
	}

	// Delete playlist (cascade will handle blocks, exercises, tags)
	return s.playlistRepo.Delete(ctx, id)
}

// ValidatePlaylistAccess checks if user can access playlist
func (s *playlistService) ValidatePlaylistAccess(ctx context.Context, playlistID int, userID uuid.UUID) error {
	playlist, err := s.playlistRepo.GetByID(ctx, playlistID)
	if err != nil {
		return err
	}

	if playlist.ID == 0 {
		return ErrPlaylistNotFound
	}

	if playlist.UserID != userID {
		return ErrUnauthorizedAccess
	}

	return nil
}
