package account

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	GetByID(ctx context.Context, id string) (Account, error)
	Create(ctx context.Context, documentNumber string) (Account, error)
}

type pgxRepository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &pgxRepository{db: db}
}

func (r *pgxRepository) GetByID(ctx context.Context, id string) (Account, error) {
	query := `SELECT id, document_number, created_at FROM account WHERE id = $1`

	row := r.db.QueryRow(ctx, query, id)

	var a Account
	err := row.Scan(
		&a.ID,
		&a.DocumentNumber,
		&a.CreatedAt,
	)
	if err != nil {
		return Account{}, fmt.Errorf("failed to get user: %w", err)
	}

	return a, nil
}

func (r *pgxRepository) Create(ctx context.Context, documentNumber string) (Account, error) {
	query := `
		INSERT INTO account (document_number) VALUES ($1)
		RETURNING id, created_at
	`

	row := r.db.QueryRow(ctx, query, documentNumber)

	var a Account
	a.DocumentNumber = documentNumber
	err := row.Scan(
		&a.ID,
		&a.CreatedAt,
	)
	if err != nil {
		return Account{}, fmt.Errorf("failed to create user: %w", err)
	}

	return a, nil
}
