package pgdb_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/dane4k/MerchShop/domain"
	"github.com/dane4k/MerchShop/infrastructure/db/pgdb"
)

func TestTransactionRepo_AddTransaction(t *testing.T) {
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := pgdb.NewTransactionRepo(db)

		mock.ExpectBegin()
		tx, err := db.Begin()
		assert.NoError(t, err)

		mock.ExpectExec("INSERT INTO transactions").
			WithArgs(100, 2, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err = repo.AddTransaction(ctx, tx, &domain.Transaction{
			Amount:     100,
			ReceiverID: 2,
			SenderID:   1,
		})
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestTransactionRepo_GetUserTransactions(t *testing.T) {
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := pgdb.NewTransactionRepo(db)
		rows := sqlmock.NewRows([]string{"username", "username", "amount"}).
			AddRow("bob", "alice", 100).
			AddRow("alice", "bob", 50)

		mock.ExpectQuery("SELECT r.username, s.username, t.amount FROM transactions t").
			WillReturnRows(rows)

		transactions, err := repo.GetUserTransactions(ctx, 1)
		assert.NoError(t, err)
		assert.Len(t, transactions, 2)
	})

	t.Run("ErrConnDone", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := pgdb.NewTransactionRepo(db)
		mock.ExpectQuery("SELECT r.username, s.username, t.amount FROM transactions t").
			WillReturnError(sql.ErrConnDone)

		_, err = repo.GetUserTransactions(ctx, 1)
		assert.Error(t, err)
	})
}
