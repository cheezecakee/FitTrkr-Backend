// Package playlist
package playlist

import (
	"time"

	"github.com/google/uuid"
)

// Custom types for better type safety
// Visibility and BlockType are used throughout the playlist domain

type Visibility string

const (
	VisibilityPrivate  Visibility = "private"
	VisibilityPublic   Visibility = "public"
	VisibilityUnlisted Visibility = "unlisted"
)

type BlockType string

const (
	BlockTypePlaylist BlockType = "playlist" // Individual exercises (replaces "regular")
	BlockTypeStandard BlockType = "standard" // Individual exercises (replaces "regular")
	BlockTypeSuperset BlockType = "superset" // 2 exercises back-to-back
	BlockTypeTriset   BlockType = "triset"   // 3 exercises back-to-back
	BlockTypeCircuit  BlockType = "circuit"  // Multiple exercises in sequence
	BlockTypeDropset  BlockType = "dropset"  // Same exercise with decreasing weight
	BlockTypeCardio   BlockType = "cardio"   // Cardio-focused block
	BlockTypeWarmup   BlockType = "warmup"   // Warm-up exercises
	BlockTypeCooldown BlockType = "cooldown" // Cool-down/stretching
)

// Tag represents playlist tags
type Tag struct {
	ID   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

// Playlist represents a workout playlist
type Playlist struct {
	ID          int       `json:"id" db:"id"`
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
	Title       string    `json:"title" db:"title"`
	Description *string   `json:"description" db:"description"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	LastWorked  bool      `json:"last_worked" db:"last_worked"`
	Visibility  Visibility    `json:"visibility" db:"visibility"` // 'private', 'public', 'unlisted'
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`

	// Joined data (not in DB)
	Tags   []Tag   `json:"tags,omitempty"`
	Blocks []Block `json:"blocks,omitempty"`
}

// Block represents an exercise block within a playlist
type Block struct {
	ID                    int    `json:"id" db:"id"`
	PlaylistID            int    `json:"playlist_id" db:"playlist_id"`
	Name                  string `json:"name" db:"name"`
	BlockType             BlockType `json:"block_type" db:"block_type"` // 'playlist', 'standard', 'superset', 'circuit', 'dropset', etc.
	BlockOrder            int    `json:"block_order" db:"block_order"`
	RestAfterBlockSeconds int    `json:"rest_after_block_seconds" db:"rest_after_block_seconds"`

	// Joined data (not in DB)
	Exercises []PlaylistExercise `json:"exercises,omitempty"`
}

// Config represents exercise configuration
type Config struct {
	ID          int      `json:"id" db:"id" example:"101"`
	Sets        *int     `json:"sets" db:"sets" example:"3"`
	RepsMin     *int     `json:"reps_min" db:"reps_min" example:"8"`
	RepsMax     *int     `json:"reps_max" db:"reps_max" example:"12"`
	Weight      *float64 `json:"weight" db:"weight" example:"50.0"`
	RestSeconds int      `json:"rest_seconds" db:"rest_seconds" example:"60"`
	Tempo       []int64  `json:"tempo" db:"tempo" example:"2,1,2,0"` // Must be exactly 4 elements: [eccentric, pause, concentric, pause]

	// Cardio fields
	DurationSeconds *int     `json:"duration_seconds" db:"duration_seconds" example:"600"`
	Distance        *float64 `json:"distance" db:"distance" example:"2.5"`
	TargetPace      *float64 `json:"target_pace" db:"target_pace" example:"5.0"`
	TargetHeartRate *int     `json:"target_heart_rate" db:"target_heart_rate" example:"140"`
	Incline         *float64 `json:"incline" db:"incline" example:"1.5"`

	Notes     *string   `json:"notes" db:"notes" example:"Keep elbows tucked"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// PlaylistExercise represents an exercise within a playlist block
type PlaylistExercise struct {
	ID            int       `json:"id" db:"id"`
	PlaylistID    int       `json:"playlist_id" db:"playlist_id"`
	ExerciseID    int       `json:"exercise_id" db:"exercise_id"`
	BlockID       int       `json:"block_id" db:"block_id"`
	ConfigID      int       `json:"config_id" db:"config_id"`
	ExerciseOrder int       `json:"exercise_order" db:"exercise_order"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`

	// Joined data (not in DB) - populated by service layer
	ExerciseName string  `json:"exercise_name,omitempty"`
	Config       *Config `json:"config,omitempty"`
}

// CreatePlaylistRequest Request/Response DTOs
type CreatePlaylistRequest struct {
	Title       string  `json:"title" validate:"required,min=1,max=100"`
	Description *string `json:"description,omitempty"`
	Visibility  string  `json:"visibility" validate:"oneof=private public unlisted"`
	GoalID      *int    `json:"goal_id,omitempty"`
	TagIDs      []int   `json:"tag_ids,omitempty"`
}

type UpdatePlaylistRequest struct {
	Title       *string `json:"title,omitempty" validate:"omitempty,min=1,max=100"`
	Description *string `json:"description,omitempty"`
	Visibility  *string `json:"visibility,omitempty" validate:"omitempty,oneof=private public unlisted"`
	GoalID      *int    `json:"goal_id,omitempty"`
	TagIDs      []int   `json:"tag_ids,omitempty"`
}

type AddExerciseToPlaylistRequest struct {
	ExerciseID int     `json:"exercise_id" validate:"required" example:"4"`
	BlockID    *int    `json:"block_id,omitempty" example:"1"` // If nil, creates new block
	BlockName  *string `json:"block_name,omitempty" example:"Push Block"`
	Config     Config  `json:"config" validate:"required"`
}

// PlaylistWithDetails includes all related data
type PlaylistWithDetails struct {
	Playlist
	TotalExercises int `json:"total_exercises"`
	TotalBlocks    int `json:"total_blocks"`
}
