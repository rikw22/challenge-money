package transaction

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, transaction *Transaction) error
}

type pgxRepository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &pgxRepository{db: db}
}

func (r pgxRepository) Create(ctx context.Context, t *Transaction) error {
	query := `
		INSERT INTO transaction (account_id, operationtype_id, amount, eventdate)
		VALUES ($1, $2, $3, $4)
		RETURNING id, eventdate
	`

	row := r.db.QueryRow(ctx, query, t.AccountId, t.OperationTypeId, t.Amount, t.EventDate)
	err := row.Scan(
		&t.ID,
		&t.EventDate,
	)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	return nil
}
