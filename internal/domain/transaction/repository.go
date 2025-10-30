package transaction

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, transaction *Transaction) error
	GetTransactionsWithNegativeBalance(ctx context.Context, accountId int) ([]Transaction, error)
	UpdateTransactionBalance(ctx context.Context, uuid pgtype.UUID, balance int) error
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

func (r *pgxRepository) GetTransactionsWithNegativeBalance(ctx context.Context, accountId int) ([]Transaction, error) {
	query := `
		SELECT id, account_id, operationtype_id, amount, balance, eventdate 
		FROM transaction WHERE account_id=$1 AND balance < 0
		ORDER BY eventdate ASC
		`
	rows, err := r.db.Query(ctx, query, accountId)
	if err != nil {
		return nil, fmt.Errorf("failed to get rows: %w", err)
	}
	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var t Transaction
		rows.Scan(
			&t.ID,
			&t.AccountId,
			&t.OperationTypeId,
			&t.Amount,
			&t.Balance,
			&t.EventDate,
		)

		transactions = append(transactions, t)
	}

	return transactions, nil
}

func (r pgxRepository) UpdateTransactionBalance(ctx context.Context, uuid pgtype.UUID, balance int) error {
	query := `UPDATE transaction SET balance=$1 WHERE id=$2`
	_, err := r.db.Query(ctx, query, balance, uuid)
	if err != nil {
		return fmt.Errorf("failed to update balance: %w", err)
	}
	return nil
}
