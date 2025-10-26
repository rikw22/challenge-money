package account

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, documentNumber string) (Account, error)
}

type pgxRepository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &pgxRepository{db: db}
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
