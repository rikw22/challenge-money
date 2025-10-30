package account

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	GetByID(ctx context.Context, id string) (Account, error)
	Create(ctx context.Context, account *Account) error
	Exist(ctx context.Context, id int) (bool, error)
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

func (r *pgxRepository) Create(ctx context.Context, a *Account) error {
	query := `
		INSERT INTO account (document_number) VALUES ($1)
		RETURNING id, created_at
	`

	row := r.db.QueryRow(ctx, query, a.DocumentNumber)

	err := row.Scan(
		&a.ID,
		&a.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *pgxRepository) Exist(ctx context.Context, id int) (bool, error) {
	query := `SELECT COUNT(id)>0 FROM account WHERE id=$1;`

	exist := false
	row := r.db.QueryRow(ctx, query, id)
	err := row.Scan(
		&exist,
	)
	if err != nil {
		return false, fmt.Errorf("failed to check if account exist: %w", err)
	}

	return exist, nil
}
