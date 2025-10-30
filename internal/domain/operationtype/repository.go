package operationtype

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Exist(ctx context.Context, id int) (bool, error)
}

type pgxRepository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &pgxRepository{db: db}
}

func (r *pgxRepository) Exist(ctx context.Context, id int) (bool, error) {
	query := `SELECT COUNT(id)>0 FROM operationtype WHERE id=$1;`

	exist := false
	row := r.db.QueryRow(ctx, query, id)
	err := row.Scan(
		&exist,
	)
	if err != nil {
		return false, fmt.Errorf("failed to check if operation type exist: %w", err)
	}

	return exist, nil
}
