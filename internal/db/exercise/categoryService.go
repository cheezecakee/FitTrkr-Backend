package exercise

import (
	"context"
)

type CategoryService interface {
	Create(ctx context.Context, category *Category) error
	GetByID(ctx context.Context, id int) (*Category, error)
	GetByName(ctx context.Context, name string) (*Category, error)
	List(ctx context.Context, offset, limit int) ([]*Category, error)
}

type DBCategoryService struct {
	repo CategoryRepo
}

func NewCategoryService(repo CategoryRepo) CategoryService {
	return &DBCategoryService{repo: repo}
}

func (s *DBCategoryService) Create(ctx context.Context, category *Category) error {
	return s.repo.Create(ctx, category)
}

func (s *DBCategoryService) GetByID(ctx context.Context, id int) (*Category, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *DBCategoryService) GetByName(ctx context.Context, name string) (*Category, error) {
	return s.repo.GetByName(ctx, name)
}

func (s *DBCategoryService) List(ctx context.Context, offset, limit int) ([]*Category, error) {
	return s.repo.List(ctx, offset, limit)
}
