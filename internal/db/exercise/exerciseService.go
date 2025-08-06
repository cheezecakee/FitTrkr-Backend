package exercise

import (
	"context"
	"errors"
)

// ExerciseService defines the interface for exercise-related operations
type ExerciseService interface {
	// Core CRUD operations
	Create(ctx context.Context, req *CreateExerciseRequest) (*Exercise, error)
	Update(ctx context.Context, updateReq *UpdateExerciseRequest, exerciseID int) (*Exercise, error)
	Delete(ctx context.Context, id int) error

	// Query operations
	GetByID(ctx context.Context, id int) (*Exercise, error)
	GetByName(ctx context.Context, name string) (*Exercise, error)
	GetByCategoryName(ctx context.Context, category string) ([]*Exercise, error)
	GetByEquipmentName(ctx context.Context, equipment string) ([]*Exercise, error)
	List(ctx context.Context, offset, limit int) ([]*Exercise, error)
	Search(ctx context.Context, query string) ([]*Exercise, error)

	// Relationship operations
	GetByMuscleGroupID(ctx context.Context, muscleGroupID int) ([]*Exercise, error)
	GetByMuscleGroupName(ctx context.Context, muscleName string) ([]*Exercise, error)
	GetByTrainingTypeID(ctx context.Context, typeID int) ([]*Exercise, error)
	GetByTrainingTypeName(ctx context.Context, typeName string) ([]*Exercise, error)

	// Enhanced operations that return full exercise details
	GetExerciseWithDetails(ctx context.Context, id int) (*Exercise, error)
	CreateWithRelations(ctx context.Context, req *CreateExerciseRequest) (*Exercise, error)
	UpdateWithRelations(ctx context.Context, req *UpdateExerciseRequest, exerciseID int) (*Exercise, error)
}

type exerciseService struct {
	repo ExerciseRepo
}

// NewExerciseService creates a new instance of ExerciseService
func NewExerciseService(repo ExerciseRepo) ExerciseService {
	return &exerciseService{
		repo: repo,
	}
}

// CreateWithRelations creates an exercise with all its relationships
func (s *exerciseService) CreateWithRelations(ctx context.Context, req *CreateExerciseRequest) (*Exercise, error) {
	// Validate inputs
	if err := s.validateCreateRequest(req); err != nil {
		return nil, err
	}

	// Check for duplicate name
	if err := s.checkDuplicateName(ctx, req.Name); err != nil {
		return nil, err
	}

	// Create the base exercise
	exercise := &Exercise{
		Name:        req.Name,
		Description: req.Description,
		Category:    &Category{ID: req.CategoryID},
		Equipment:   &Equipment{ID: req.EquipmentID},
	}

	if err := s.repo.Create(ctx, exercise); err != nil {
		return nil, err
	}

	// Add muscle group relationships
	if len(req.MuscleGroupIDs) > 0 {
		if err := s.repo.AddMuscleGroups(ctx, exercise.ID, req.MuscleGroupIDs); err != nil {
			return nil, err
		}
	}

	// Add training type relationships
	if len(req.TypeIDs) > 0 {
		if err := s.repo.AddExerciseTypes(ctx, exercise.ID, req.TypeIDs); err != nil {
			return nil, err
		}
	}

	// Return the exercise with all details
	return s.GetExerciseWithDetails(ctx, exercise.ID)
}

// UpdateWithRelations updates an exercise and all its relationships
func (s *exerciseService) UpdateWithRelations(ctx context.Context, req *UpdateExerciseRequest, exerciseID int) (*Exercise, error) {
	// Validate inputs
	if err := s.validateUpdateRequest(req, exerciseID); err != nil {
		return nil, err
	}

	// Update the base exercise
	exercise := &Exercise{
		ID:          exerciseID,
		Name:        req.Name,
		Description: req.Description,
		Category:    &Category{ID: req.CategoryID},
		Equipment:   &Equipment{ID: req.EquipmentID},
	}

	if err := s.repo.Update(ctx, exercise); err != nil {
		return nil, err
	}

	// Replace muscle group relationships
	if err := s.repo.RemoveAllMuscleGroups(ctx, exerciseID); err != nil {
		return nil, err
	}
	if len(req.MuscleGroupIDs) > 0 {
		if err := s.repo.AddMuscleGroups(ctx, exerciseID, req.MuscleGroupIDs); err != nil {
			return nil, err
		}
	}

	// Replace training type relationships
	if err := s.repo.RemoveAllExerciseTypes(ctx, exerciseID); err != nil {
		return nil, err
	}
	if len(req.TypeIDs) > 0 {
		if err := s.repo.AddExerciseTypes(ctx, exerciseID, req.TypeIDs); err != nil {
			return nil, err
		}
	}

	// Return the updated exercise with all details
	return s.GetExerciseWithDetails(ctx, exerciseID)
}

// GetExerciseWithDetails returns an exercise with all related data populated
func (s *exerciseService) GetExerciseWithDetails(ctx context.Context, id int) (*Exercise, error) {
	if id == 0 {
		return nil, errors.New("valid exercise ID is required")
	}

	// Get the base exercise
	exercise, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Get muscle groups
	muscleGroups, err := s.repo.GetMuscleGroups(ctx, id)
	if err == nil && len(muscleGroups) > 0 {
		exercise.MuscleGroups = make([]MuscleGroup, len(muscleGroups))
		for i, mg := range muscleGroups {
			exercise.MuscleGroups[i] = *mg
		}
	}

	// Get training types
	trainingTypes, err := s.repo.GetExerciseTypes(ctx, id)
	if err == nil && len(trainingTypes) > 0 {
		exercise.Types = make([]TrainingType, len(trainingTypes))
		for i, tt := range trainingTypes {
			exercise.Types[i] = *tt
		}
	}

	return exercise, nil
}

// Core CRUD operations (simplified versions)
func (s *exerciseService) Create(ctx context.Context, req *CreateExerciseRequest) (*Exercise, error) {
	return s.CreateWithRelations(ctx, req)
}

func (s *exerciseService) Update(ctx context.Context, req *UpdateExerciseRequest, exerciseID int) (*Exercise, error) {
	return s.UpdateWithRelations(ctx, req, exerciseID)
}

func (s *exerciseService) Delete(ctx context.Context, id int) error {
	if id == 0 {
		return errors.New("valid exercise ID is required")
	}
	return s.repo.Delete(ctx, id)
}

func (s *exerciseService) GetByID(ctx context.Context, id int) (*Exercise, error) {
	if id == 0 {
		return nil, errors.New("valid exercise ID is required")
	}
	return s.repo.GetByID(ctx, id)
}

func (s *exerciseService) GetByName(ctx context.Context, name string) (*Exercise, error) {
	if name == "" {
		return nil, errors.New("exercise name is required")
	}
	return s.repo.GetByName(ctx, name)
}

func (s *exerciseService) GetByCategoryName(ctx context.Context, category string) ([]*Exercise, error) {
	if category == "" {
		return nil, errors.New("category name is required")
	}
	return s.repo.GetByCategoryID(ctx, category)
}

func (s *exerciseService) GetByEquipmentName(ctx context.Context, equipment string) ([]*Exercise, error) {
	if equipment == "" {
		return nil, errors.New("equipment name is required")
	}
	return s.repo.GetByEquipmentName(ctx, equipment)
}

func (s *exerciseService) List(ctx context.Context, offset, limit int) ([]*Exercise, error) {
	if limit <= 0 {
		return nil, errors.New("limit must be greater than 0")
	}
	return s.repo.List(ctx, offset, limit)
}

func (s *exerciseService) Search(ctx context.Context, query string) ([]*Exercise, error) {
	if query == "" {
		return nil, errors.New("search query is required")
	}
	return s.repo.Search(ctx, query)
}

// Relationship query operations
func (s *exerciseService) GetByMuscleGroupID(ctx context.Context, muscleGroupID int) ([]*Exercise, error) {
	if muscleGroupID == 0 {
		return nil, errors.New("valid muscle group ID is required")
	}
	return s.repo.GetExercisesByMuscle(ctx, muscleGroupID)
}

func (s *exerciseService) GetByMuscleGroupName(ctx context.Context, muscleName string) ([]*Exercise, error) {
	if muscleName == "" {
		return nil, errors.New("muscle group name is required")
	}
	return s.repo.GetExercisesByMuscleName(ctx, muscleName)
}

func (s *exerciseService) GetByTrainingTypeID(ctx context.Context, typeID int) ([]*Exercise, error) {
	if typeID == 0 {
		return nil, errors.New("valid training type ID is required")
	}
	return s.repo.GetExercisesByType(ctx, typeID)
}

func (s *exerciseService) GetByTrainingTypeName(ctx context.Context, typeName string) ([]*Exercise, error) {
	if typeName == "" {
		return nil, errors.New("training type name is required")
	}
	return s.repo.GetExercisesByTypeName(ctx, typeName)
}

// Helper validation methods
func (s *exerciseService) validateCreateRequest(req *CreateExerciseRequest) error {
	if req.Name == "" {
		return errors.New("exercise name is required")
	}
	if len(req.Name) > 100 {
		return errors.New("exercise name must not exceed 100 characters")
	}
	if req.Description == "" {
		return errors.New("exercise description is required")
	}
	if req.CategoryID <= 0 {
		return errors.New("valid category ID is required")
	}
	if req.EquipmentID <= 0 {
		return errors.New("valid equipment ID is required")
	}
	return nil
}

func (s *exerciseService) validateUpdateRequest(req *UpdateExerciseRequest, exerciseID int) error {
	if exerciseID == 0 {
		return errors.New("valid exercise ID is required")
	}
	return s.validateCreateRequest(&CreateExerciseRequest{
		Name:           req.Name,
		Description:    req.Description,
		CategoryID:     req.CategoryID,
		EquipmentID:    req.EquipmentID,
		TypeIDs:        req.TypeIDs,
		MuscleGroupIDs: req.MuscleGroupIDs,
	})
}

func (s *exerciseService) checkDuplicateName(ctx context.Context, name string) error {
	existing, err := s.repo.GetByName(ctx, name)
	if err == nil && existing != nil && existing.ID > 0 {
		return errors.New("exercise with this name already exists")
	}
	return nil
}
