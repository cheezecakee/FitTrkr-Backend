package exercise

import (
	"context"
	"database/sql"

	"github.com/cheezecakee/fitrkr/internal/utils/transaction"
)

type EquipmentRepo interface {
	Create(ctx context.Context, equipment *Equipment) error
	GetByID(ctx context.Context, id int) (*Equipment, error)
	GetByName(ctx context.Context, name string) (*Equipment, error)
	List(ctx context.Context, offset, limit int) ([]*Equipment, error)
}

type DBEquipmentRepo struct {
	tx transaction.BaseRepository
}

func NewEquipmentRepo(db *sql.DB) EquipmentRepo {
	return &DBEquipmentRepo{
		tx: transaction.NewBaseRepository(db),
	}
}

const createEquipment = `INSERT INTO equipment (name) VALUES ($1) RETURNING id`

func (r *DBEquipmentRepo) Create(ctx context.Context, equipment *Equipment) error {
	return r.tx.WithTransaction(ctx, func(tx *sql.Tx) error {
		return tx.QueryRowContext(ctx, createEquipment, equipment.Name).Scan(&equipment.ID)
	})
}

const getEquipmentByID = `SELECT id, name FROM equipment WHERE id = $1 LIMIT 1`

func (r *DBEquipmentRepo) GetByID(ctx context.Context, id int) (*Equipment, error) {
	equipment := &Equipment{}
	err := r.tx.DB().QueryRowContext(ctx, getEquipmentByID, id).Scan(&equipment.ID, &equipment.Name)
	return equipment, err
}

const getEquipmentByName = `SELECT id, name FROM equipment WHERE name = $1 LIMIT 1`

func (r *DBEquipmentRepo) GetByName(ctx context.Context, name string) (*Equipment, error) {
	equipment := &Equipment{}
	err := r.tx.DB().QueryRowContext(ctx, getEquipmentByName, name).Scan(&equipment.ID, &equipment.Name)
	return equipment, err
}

const listEquipment = `SELECT id, name FROM equipment OFFSET $1 LIMIT $2`

func (r *DBEquipmentRepo) List(ctx context.Context, offset, limit int) ([]*Equipment, error) {
	rows, err := r.tx.DB().QueryContext(ctx, listEquipment, offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var equipmentList []*Equipment
	for rows.Next() {
		equipment := &Equipment{}
		if err := rows.Scan(&equipment.ID, &equipment.Name); err != nil {
			return nil, err
		}
		equipmentList = append(equipmentList, equipment)
	}
	return equipmentList, rows.Err()
}
