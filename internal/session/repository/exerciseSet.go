package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"

	m "github/cheezecakee/fitrkr/internal/session/models"
)

type DBExSetRepo struct {
	db *sql.DB
}

func NewSession(db *sql.DB) ExSetRepo {
	return &DBExSetRepo{
		db: db,
	}
}

const createExSet = ``

func (r *DBExSetRepo) Create(ctx context.Context, set *m.ExSet) error {}

const createBatchExSet = ``

func (r *DBExSetRepo) CreateBatch(ctx context.Context, sets []*m.ExSet) error {}

const getExSetByID = ``

func (r *DBExSetRepo) GetByID(ctx context.Context, id uint) (*m.ExSet, error) {}

const getExSetBySessionExID = ``

func (r *DBExSetRepo) GetBySessionExID(ctx context.Context, sessionExID uuid.UUID) ([]*m.ExSet, error) {
}

const updateExSet = ``

func (r *DBExSetRepo) Update(ctx context.Context, set *m.ExSet) error {}

const deleteExSet = ``

func (r *DBExSetRepo) Delete(ctx context.Context, id uint) error {}
