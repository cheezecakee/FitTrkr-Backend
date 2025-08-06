package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/cheezecakee/fitrkr/internal/db/playlist"
)

// PlaylistHandler handles HTTP requests for playlist operations
type PlaylistHandler struct {
	playlistSvc playlist.PlaylistService
}

// NewPlaylistHandler creates a new playlist handler
func NewPlaylistHandler(playlistSvc playlist.PlaylistService) *PlaylistHandler {
	return &PlaylistHandler{
		playlistSvc: playlistSvc,
	}
}

// CreatePlaylist godoc
// @Summary Create a new playlist
// @Description Create a new workout playlist for the authenticated user
// @Tags playlists
// @Accept json
// @Produce json
// @Param request body playlist.CreatePlaylistRequest true "Playlist creation request"
// @Success 201 {object} playlist.Playlist "Created playlist"
// @Failure 400 {object} errors.ErrorResponse "Bad request"
// @Failure 401 {object} errors.ErrorResponse "Unauthorized"
// @Failure 409 {object} errors.ErrorResponse "Playlist already exists"
// @Failure 500 {object} errors.ErrorResponse "Internal server error"
// @Router /api/v1/playlists [post]
// @Security BearerAuth
func (h *PlaylistHandler) CreatePlaylist(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(UserIDKey).(uuid.UUID)
	if !ok {
		ClientError(w, http.StatusUnauthorized)
		return
	}

	var req playlist.CreatePlaylistRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Basic validation
	if req.Title == "" {
		ErrorResponse(w, http.StatusBadRequest, "Title is required")
		return
	}

	createdPlaylist, err := h.playlistSvc.CreatePlaylist(r.Context(), userID, req)
	if err != nil {
		switch err {
		case playlist.ErrPlaylistExists:
			ErrorResponse(w, http.StatusConflict, "Playlist with this title already exists")
		default:
			ServerError(w, err)
		}
		return
	}

	Response(w, http.StatusCreated, createdPlaylist)
}

// GetPlaylist godoc
// @Summary Get a playlist by ID
// @Description Get playlist details by ID
// @Tags playlists
// @Produce json
// @Param id path int true "Playlist ID"
// @Success 200 {object} playlist.Playlist "Playlist details"
// @Failure 400 {object} errors.ErrorResponse "Bad request"
// @Failure 401 {object} errors.ErrorResponse "Unauthorized"
// @Failure 403 {object} errors.ErrorResponse "Forbidden"
// @Failure 404 {object} errors.ErrorResponse "Playlist not found"
// @Failure 500 {object} errors.ErrorResponse "Internal server error"
// @Router /api/v1/playlists/{id} [get]
// @Security BearerAuth
func (h *PlaylistHandler) GetPlaylist(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(UserIDKey).(uuid.UUID)
	if !ok {
		ClientError(w, http.StatusUnauthorized)
		return
	}

	playlistID, err := h.extractPlaylistID(r)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid playlist ID")
		return
	}
	playlistData, err := h.playlistSvc.GetPlaylistByID(r.Context(), playlistID, userID)
	if err != nil {
		switch err {
		case playlist.ErrPlaylistNotFound:
			ErrorResponse(w, http.StatusNotFound, "Playlist not found")
		case playlist.ErrUnauthorizedAccess:
			ErrorResponse(w, http.StatusForbidden, "Access denied")
		default:
			ServerError(w, err)
		}
		return
	}

	Response(w, http.StatusOK, playlistData)
}

// GetUserPlaylists godoc
// @Summary Get user's playlists
// @Description Get all playlists for the authenticated user
// @Tags playlists
// @Produce json
// @Success 200 {array} playlist.PlaylistWithDetails "User's playlists"
// @Failure 401 {object} errors.ErrorResponse "Unauthorized"
// @Failure 500 {object} errors.ErrorResponse "Internal server error"
// @Router /api/v1/playlists [get]
// @Security BearerAuth
func (h *PlaylistHandler) GetUserPlaylists(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(UserIDKey).(uuid.UUID)
	if !ok {
		ClientError(w, http.StatusUnauthorized)
		return
	}

	playlists, err := h.playlistSvc.GetUserPlaylists(r.Context(), userID)
	if err != nil {
		ServerError(w, err)
		return
	}

	Response(w, http.StatusOK, playlists)
}

// GetPlaylistForSession godoc
// @Summary Get playlist for workout session
// @Description Get complete playlist data including exercises and configs for starting a workout session
// @Tags playlists
// @Produce json
// @Param id path int true "Playlist ID"
// @Success 200 {object} playlist.Playlist "Complete playlist data"
// @Failure 400 {object} errors.ErrorResponse "Bad request"
// @Failure 401 {object} errors.ErrorResponse "Unauthorized"
// @Failure 403 {object} errors.ErrorResponse "Forbidden"
// @Failure 404 {object} errors.ErrorResponse "Playlist not found"
// @Failure 500 {object} errors.ErrorResponse "Internal server error"
// @Router /api/v1/playlists/{id}/session [get]
// @Security BearerAuth
func (h *PlaylistHandler) GetPlaylistForSession(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(UserIDKey).(uuid.UUID)
	if !ok {
		ClientError(w, http.StatusUnauthorized)
		return
	}

	playlistID, err := h.extractPlaylistID(r)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid playlist ID")
		return
	}

	sessionPlaylist, err := h.playlistSvc.GetPlaylistForSession(r.Context(), playlistID, userID)
	if err != nil {
		switch err {
		case playlist.ErrPlaylistNotFound:
			ErrorResponse(w, http.StatusNotFound, "Playlist not found")
		case playlist.ErrUnauthorizedAccess:
			ErrorResponse(w, http.StatusForbidden, "Access denied")
		default:
			ServerError(w, err)
		}
		return
	}

	Response(w, http.StatusOK, sessionPlaylist)
}

// UpdatePlaylist godoc
// @Summary Update a playlist
// @Description Update playlist details
// @Tags playlists
// @Accept json
// @Produce json
// @Param id path int true "Playlist ID"
// @Param request body playlist.UpdatePlaylistRequest true "Playlist update request"
// @Success 200 {object} playlist.Playlist "Updated playlist"
// @Failure 400 {object} errors.ErrorResponse "Bad request"
// @Failure 401 {object} errors.ErrorResponse "Unauthorized"
// @Failure 403 {object} errors.ErrorResponse "Forbidden"
// @Failure 404 {object} errors.ErrorResponse "Playlist not found"
// @Failure 409 {object} errors.ErrorResponse "Playlist title already exists"
// @Failure 500 {object} errors.ErrorResponse "Internal server error"
// @Router /api/v1/playlists/{id} [put]
// @Security BearerAuth
func (h *PlaylistHandler) UpdatePlaylist(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(UserIDKey).(uuid.UUID)
	if !ok {
		ClientError(w, http.StatusUnauthorized)
		return
	}

	playlistID, err := h.extractPlaylistID(r)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid playlist ID")
		return
	}

	var req playlist.UpdatePlaylistRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	updatedPlaylist, err := h.playlistSvc.UpdatePlaylist(r.Context(), playlistID, userID, req)
	if err != nil {
		switch err {
		case playlist.ErrPlaylistNotFound:
			ErrorResponse(w, http.StatusNotFound, "Playlist not found")
		case playlist.ErrUnauthorizedAccess:
			ErrorResponse(w, http.StatusForbidden, "Access denied")
		case playlist.ErrPlaylistExists:
			ErrorResponse(w, http.StatusConflict, "Playlist with this title already exists")
		default:
			ServerError(w, err)
		}
		return
	}

	Response(w, http.StatusOK, updatedPlaylist)
}

// DeletePlaylist godoc
// @Summary Delete a playlist
// @Description Delete a playlist and all associated data
// @Tags playlists
// @Param id path int true "Playlist ID"
// @Success 204 "Playlist deleted successfully"
// @Failure 400 {object} errors.ErrorResponse "Bad request"
// @Failure 401 {object} errors.ErrorResponse "Unauthorized"
// @Failure 403 {object} errors.ErrorResponse "Forbidden"
// @Failure 404 {object} errors.ErrorResponse "Playlist not found"
// @Failure 500 {object} errors.ErrorResponse "Internal server error"
// @Router /api/v1/playlists/{id} [delete]
// @Security BearerAuth
func (h *PlaylistHandler) DeletePlaylist(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(UserIDKey).(uuid.UUID)
	if !ok {
		ClientError(w, http.StatusUnauthorized)
		return
	}

	playlistID, err := h.extractPlaylistID(r)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid playlist ID")
		return
	}

	err = h.playlistSvc.DeletePlaylist(r.Context(), playlistID, userID)
	if err != nil {
		switch err {
		case playlist.ErrPlaylistNotFound:
			ErrorResponse(w, http.StatusNotFound, "Playlist not found")
		case playlist.ErrUnauthorizedAccess:
			ErrorResponse(w, http.StatusForbidden, "Access denied")
		default:
			ServerError(w, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// AddExerciseToPlaylist godoc
// @Summary Add exercise to playlist
// @Description Add an exercise to a playlist with configuration
// @Tags playlists
// @Accept json
// @Produce json
// @Param id path int true "Playlist ID"
// @Param request body playlist.AddExerciseToPlaylistRequest true "Add exercise request"
// @Success 201 {object} playlist.PlaylistExercise "Added exercise"
// @Failure 400 {object} errors.ErrorResponse "Bad request"
// @Failure 401 {object} errors.ErrorResponse "Unauthorized"
// @Failure 403 {object} errors.ErrorResponse "Forbidden"
// @Failure 404 {object} errors.ErrorResponse "Playlist or block not found"
// @Failure 500 {object} errors.ErrorResponse "Internal server error"
// @Router /api/v1/playlists/{id}/exercises [post]
// @Security BearerAuth
func (h *PlaylistHandler) AddExerciseToPlaylist(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(UserIDKey).(uuid.UUID)
	if !ok {
		ClientError(w, http.StatusUnauthorized)
		return
	}

	playlistID, err := h.extractPlaylistID(r)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid playlist ID")
		return
	}

	var req playlist.AddExerciseToPlaylistRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Basic validation
	if req.ExerciseID == 0 {
		ErrorResponse(w, http.StatusBadRequest, "Exercise ID is required")
		return
	}

	addedExercise, err := h.playlistSvc.AddExerciseToPlaylist(r.Context(), playlistID, userID, req)
	if err != nil {
		switch err {
		case playlist.ErrPlaylistNotFound:
			ErrorResponse(w, http.StatusNotFound, "Playlist not found")
		case playlist.ErrUnauthorizedAccess:
			ErrorResponse(w, http.StatusForbidden, "Access denied")
		case playlist.ErrBlockNotFound:
			ErrorResponse(w, http.StatusNotFound, "Block not found")
		default:
			ServerError(w, err)
		}
		return
	}

	Response(w, http.StatusCreated, addedExercise)
}

// RemoveExerciseFromPlaylist godoc
// @Summary Remove exercise from playlist
// @Description Remove an exercise from a playlist
// @Tags playlists
// @Param id path int true "Playlist Exercise ID"
// @Success 204 "Exercise removed successfully"
// @Failure 400 {object} errors.ErrorResponse "Bad request"
// @Failure 401 {object} errors.ErrorResponse "Unauthorized"
// @Failure 403 {object} errors.ErrorResponse "Forbidden"
// @Failure 404 {object} errors.ErrorResponse "Exercise not found"
// @Failure 500 {object} errors.ErrorResponse "Internal server error"
// @Router /api/v1/playlists/exercises/{id} [delete]
// @Security BearerAuth
func (h *PlaylistHandler) RemoveExerciseFromPlaylist(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(UserIDKey).(uuid.UUID)
	if !ok {
		ClientError(w, http.StatusUnauthorized)
		return
	}

	exerciseIDStr := chi.URLParam(r, "id")
	exerciseID, err := strconv.Atoi(exerciseIDStr)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid exercise ID")
		return
	}

	err = h.playlistSvc.RemoveExerciseFromPlaylist(r.Context(), exerciseID, userID)
	if err != nil {
		switch err {
		case playlist.ErrPlaylistNotFound:
			ErrorResponse(w, http.StatusNotFound, "Exercise not found in playlist")
		case playlist.ErrUnauthorizedAccess:
			ErrorResponse(w, http.StatusForbidden, "Access denied")
		default:
			ServerError(w, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// CreateExerciseBlock godoc
// @Summary Create exercise block
// @Description Create a new exercise block in a playlist
// @Tags playlists
// @Accept json
// @Produce json
// @Param id path int true "Playlist ID"
// @Param request body CreateBlockRequest true "Create block request"
// @Success 201 {object} playlist.Block "Created block"
// @Failure 400 {object} errors.ErrorResponse "Bad request"
// @Failure 401 {object} errors.ErrorResponse "Unauthorized"
// @Failure 403 {object} errors.ErrorResponse "Forbidden"
// @Failure 404 {object} errors.ErrorResponse "Playlist not found"
// @Failure 500 {object} errors.ErrorResponse "Internal server error"
// @Router /api/v1/playlists/{id}/blocks [post]
// @Security BearerAuth
func (h *PlaylistHandler) CreateExerciseBlock(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(UserIDKey).(uuid.UUID)
	if !ok {
		ClientError(w, http.StatusUnauthorized)
		return
	}

	playlistID, err := h.extractPlaylistID(r)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid playlist ID")
		return
	}

	var req CreateBlockRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Basic validation
	if req.Name == "" {
		ErrorResponse(w, http.StatusBadRequest, "Block name is required")
		return
	}
	if req.BlockType == "" {
		req.BlockType = string(playlist.BlockTypePlaylist)
	}

	createdBlock, err := h.playlistSvc.CreateBlock(r.Context(), playlistID, userID, req.Name, req.BlockType)
	if err != nil {
		switch err {
		case playlist.ErrPlaylistNotFound:
			ErrorResponse(w, http.StatusNotFound, "Playlist not found")
		case playlist.ErrUnauthorizedAccess:
			ErrorResponse(w, http.StatusForbidden, "Access denied")
		case playlist.ErrInvalidBlockType:
			ErrorResponse(w, http.StatusBadRequest, "Invalid block type")
		default:
			ServerError(w, err)
		}
		return
	}

	Response(w, http.StatusCreated, createdBlock)
}

// GetTags godoc
// @Summary Get all tags
// @Description Get all available playlist tags
// @Tags playlists
// @Produce json
// @Success 200 {array} playlist.Tag "Available tags"
// @Failure 500 {object} errors.ErrorResponse "Internal server error"
// @Router /api/v1/playlists/tags [get]
// @Security BearerAuth
func (h *PlaylistHandler) GetTags(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(UserIDKey).(uuid.UUID)
	if !ok {
		ClientError(w, http.StatusUnauthorized)
		return
	}
	tags, err := h.playlistSvc.GetAllTags(r.Context(), userID)
	if err != nil {
		ServerError(w, err)
		return
	}

	Response(w, http.StatusOK, tags)
}

func (h *PlaylistHandler) extractPlaylistID(r *http.Request) (int, error) {
	playlistIDStr := chi.URLParam(r, "id")
	return strconv.Atoi(playlistIDStr)
}

// CreateBlockRequest represents the request structure for creating exercise blocks
type CreateBlockRequest struct {
	Name      string `json:"name" validate:"required"`
	BlockType string `json:"block_type,omitempty"`
}
