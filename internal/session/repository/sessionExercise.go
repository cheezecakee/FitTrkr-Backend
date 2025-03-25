package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"

	m "github/cheezecakee/fitrkr/internal/session/models"
)

type DBSessionExRepo struct {
	db *sql.DB
}

func NewSessionEx(db *sql.DB) SessionExRepo {
	return &DBSessionExRepo{
		db: db,
	}
}

const createSessionEx = ``

func (r *DBSessionExRepo) Create(ctx context.Context, sessionEx *m.SessionEx) (*m.SessionEx, error) {}

const getSessionExByID = ``

func (r *DBSessionExRepo) GetByID(ctx context.Context, id uuid.UUID) (*m.SessionEx, error) {}

const getSessionExBySessionID = ``

func (r *DBSessionExRepo) GetBysessionID(ctx context.Context, sessionID uuid.UUID) ([]*m.SessionEx, error) {
}

const updateSessionEx = ``

func (r *DBSessionExRepo) Update(ctx context.Context, sessionEx *m.SessionEx) error {}

const deleteSessionEx = ``

func (r *DBSessionExRepo) Delete(ctx context.Context, id uuid.UUID) error {}

const deleteSessionExBySessionID = ``

func (r *DBSessionExRepo) DeleteBysessionID(ctx context.Context, id uuid.UUID) error {}
