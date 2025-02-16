package pgdb

import (
	"context"
	"database/sql"

	"github.com/Masterminds/squirrel"

	"github.com/dane4k/MerchShop/domain"
	"github.com/dane4k/MerchShop/internal/consts"
	"github.com/dane4k/MerchShop/presentations/dto/response"
)

type TransactionRepo struct {
	SQL squirrel.StatementBuilderType
}

func NewTransactionRepo(db *sql.DB) *TransactionRepo {
	return &TransactionRepo{
		SQL: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).RunWith(db),
	}
}

func (tr *TransactionRepo) AddTransaction(ctx context.Context, tx *sql.Tx, transaction *domain.Transaction) error {
	if transaction == nil {
		return RespondWithError(consts.ErrInternalServerError, "transaction object is nil")
	}

	if err := checkTx(tx); err != nil {
		return RespondWithError(err, "invalid transaction (sql tx)")
	}

	_, err := tr.SQL.Insert(consts.TransactionsTable).
		Columns(consts.ColumnAmount, consts.ColumnReceiverID, consts.ColumnSenderID).
		Values(transaction.Amount, transaction.ReceiverID, transaction.SenderID).
		RunWith(tx).
		ExecContext(ctx)
	if err != nil {
		return RespondWithError(err, "error adding transaction")
	}

	return nil
}

func (tr *TransactionRepo) GetUserTransactions(ctx context.Context, userID int) ([]*response.Transaction, error) {
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

	var transactions []*response.Transaction
	for rows.Next() {
		var transaction response.Transaction
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
