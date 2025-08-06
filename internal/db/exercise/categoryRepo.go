package exercise

import (
	"context"
	"database/sql"

	"github.com/cheezecakee/fitrkr/internal/utils/transaction"
)

type CategoryRepo interface {
	Create(ctx context.Context, category *Category) error
	GetByID(ctx context.Context, id int) (*Category, error)
	GetByName(ctx context.Context, name string) (*Category, error)
	List(ctx context.Context, offset, limit int) ([]*Category, error)
}

type DBCategoryRepo struct {
	tx transaction.BaseRepository
}

func NewCategoryRepo(db *sql.DB) CategoryRepo {
	return &DBCategoryRepo{
		tx: transaction.NewBaseRepository(db),
	}
}

const createCategory = `INSERT INTO exercise_categories (name) VALUES ($1) RETURNING id`

func (r *DBCategoryRepo) Create(ctx context.Context, category *Category) error {
	return r.tx.WithTransaction(ctx, func(tx *sql.Tx) error {
		return tx.QueryRowContext(ctx, createCategory, category.Name).Scan(&category.ID)
	})
}

const getCategoryByID = `SELECT id, name FROM exercise_categories WHERE id = $1 LIMIT 1`

func (r *DBCategoryRepo) GetByID(ctx context.Context, id int) (*Category, error) {
	category := &Category{}
	err := r.tx.DB().QueryRowContext(ctx, getCategoryByID, id).Scan(&category.ID, &category.Name)
	return category, err
}

const getCategoryByName = `SELECT id, name FROM exercise_categories WHERE name = $1 LIMIT 1`

func (r *DBCategoryRepo) GetByName(ctx context.Context, name string) (*Category, error) {
	category := &Category{}
	err := r.tx.DB().QueryRowContext(ctx, getCategoryByName, name).Scan(&category.ID, &category.Name)
	return category, err
}

const listCategories = `SELECT id, name FROM exercise_categories OFFSET $1 LIMIT $2`

func (r *DBCategoryRepo) List(ctx context.Context, offset, limit int) ([]*Category, error) {
	rows, err := r.tx.DB().QueryContext(ctx, listCategories, offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*Category
	for rows.Next() {
		category := &Category{}
		if err := rows.Scan(&category.ID, &category.Name); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, rows.Err()
}
