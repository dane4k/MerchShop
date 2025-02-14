package pgdb

import (
	"MerchShop/internal/dto"
	"MerchShop/internal/entity"
	"context"
	"database/sql"
	"github.com/Masterminds/squirrel"
)

type TransactionRepo interface {
	GetUserTransactions(ctx context.Context, userID int) ([]*dto.Transaction, error)
	addTransaction(ctx context.Context, tx *sql.Tx, transaction *entity.Transaction) error
}

type transactionRepo struct {
	SQL squirrel.StatementBuilderType
}

func NewTransactionRepo(db *sql.DB) TransactionRepo {
	return &transactionRepo{
		SQL: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).RunWith(db),
	}
}

func (tr *transactionRepo) addTransaction(ctx context.Context, tx *sql.Tx, transaction *entity.Transaction) error {
	if transaction == nil {
		return RespondWithError(ErrInternalServerError, "transaction object is nil")
	}

	if err := checkTx(tx); err != nil {
		return RespondWithError(err, "invalid transaction (sql tx)")
	}

	_, err := tr.SQL.Insert(TransactionsTable).
		Columns(ColumnAmount, ColumnReceiverID, ColumnSenderID).
		Values(transaction.Amount, transaction.ReceiverID, transaction.SenderID).
		RunWith(tx).
		ExecContext(ctx)
	if err != nil {
		return RespondWithError(err, "error adding transaction")
	}

	return nil
}

func (tr *transactionRepo) GetUserTransactions(ctx context.Context, userID int) ([]*dto.Transaction, error) {
	rows, err := tr.SQL.Select("r.username", "s.username", "t.amount").
		From("transactions t").
		Join("users r ON t.receiver_id = r.id").
		Join("users s ON t.sender_id = s.id").
		Where(squirrel.Or{
			squirrel.Eq{"t.receiver_id": userID},
			squirrel.Eq{"t.sender_id": userID},
		}).QueryContext(ctx)
	if err != nil {
		return nil, RespondWithError(err, "error getting user transactions")
	}
	defer rows.Close()

	var transactions []*dto.Transaction
	for rows.Next() {
		var transaction dto.Transaction
		err = rows.Scan(&transaction.ReceiverUsername, &transaction.SenderUsername, &transaction.Amount)
		if err != nil {
			return nil, RespondWithError(err, "error getting user transactions")
		}
		transactions = append(transactions, &transaction)
	}
	if err = rows.Err(); err != nil {
		return nil, RespondWithError(err, "error getting user transactions")
	}

	return transactions, nil
}
