// Package exercise
package exercise

import (
	"time"
)

type Exercise struct {
	ID           int            `db:"id" json:"id"`
	Name         string         `db:"name" json:"name"`
	Description  string         `db:"description" json:"description"`
	Category     *Category      `json:"category,omitempty"`
	Equipment    *Equipment     `json:"equipment,omitempty"`
	Types        []TrainingType `json:"training_types,omitempty"`
	MuscleGroups []MuscleGroup  `json:"muscleGroups,omitempty"`
	CreatedAt    time.Time      `db:"created_at" json:"createdAt"`
	UpdatedAt    time.Time      `db:"updated_at" json:"updatedAt"`
}

type ExerciseResponse struct {
	ID           int            `db:"id" json:"id"`
	Name         string         `db:"name" json:"name"`
	Description  string         `db:"description" json:"description"`
	CategoryID   int            `db:"category_id" json:"categoryID"`
	EquipmentID  int            `db:"equipment_id" json:"equipmentID"`
	Types        []TrainingType `json:"training_types,omitempty"`
	MuscleGroups []MuscleGroup  `json:"muscleGroups,omitempty"`
	CreatedAt    time.Time      `db:"created_at" json:"createdAt"`
	UpdatedAt    time.Time      `db:"updated_at" json:"updatedAt"`
}

// UpdateExerciseRequest swagger:model UpdateExerciseRequest
// UpdateExerciseRequest is used for updating an exercise, including types and muscle groups.
type UpdateExerciseRequest struct {
	Name           string `json:"name"`
	Description    string `json:"description"`
	CategoryID     int    `json:"categoryID"`
	EquipmentID    int    `json:"equipmentID"`
	TypeIDs        []int  `json:"typeIDs"`
	MuscleGroupIDs []int  `json:"muscleGroupIDs"`
}

// CreateExerciseRequest swagger:model CreateExerciseRequest
// CreateExerciseRequest is used for creating an exercise, including types and muscle groups.
type CreateExerciseRequest struct {
	Name           string `json:"name"`
	Description    string `json:"description"`
	CategoryID     int    `json:"categoryID"`
	EquipmentID    int    `json:"equipmentID"`
	TypeIDs        []int  `json:"typeIDs"`
	MuscleGroupIDs []int  `json:"muscleGroupIDs"`
}

type Category struct {
	ID   int
	Name string
}

type Equipment struct {
	ID   int
	Name string
}

type TrainingType struct {
	ID   int
	Name string
}

type MuscleGroup struct {
	ID   int
	Name string
}
