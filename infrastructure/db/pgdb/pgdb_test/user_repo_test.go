package pgdb_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/dane4k/MerchShop/domain"
	"github.com/dane4k/MerchShop/infrastructure/db/pgdb"
	"github.com/dane4k/MerchShop/internal/consts"
)

func TestUserRepo_GetUserBalance(t *testing.T) {
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := pgdb.NewUserRepo(db, pgdb.NewTransactionRepo(db), pgdb.NewInventoryRepo(db))
		rows := sqlmock.NewRows([]string{"coins"}).AddRow(1000)

		mock.ExpectQuery("SELECT coins FROM users").
			WithArgs(1).
			WillReturnRows(rows)

		balance, err := repo.GetUserBalance(ctx, 1)
		assert.NoError(t, err)
		assert.Equal(t, 1000, balance)
	})

	t.Run("User not found", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := pgdb.NewUserRepo(db, pgdb.NewTransactionRepo(db), pgdb.NewInventoryRepo(db))
		mock.ExpectQuery("SELECT coins FROM users").
			WithArgs(1).
			WillReturnError(sql.ErrNoRows)

		_, err = repo.GetUserBalance(ctx, 1)
		assert.Error(t, err)
	})
}

func TestUserRepo_SendCoins(t *testing.T) {
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		txRepo := pgdb.NewTransactionRepo(db)
		invRepo := pgdb.NewInventoryRepo(db)
		repo := pgdb.NewUserRepo(db, txRepo, invRepo)

		mock.ExpectBegin()
		mock.ExpectQuery("SELECT coins FROM users WHERE id = ?").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"coins"}).AddRow(1000))

		mock.ExpectExec("UPDATE users SET coins = CASE").
			WithArgs(1, 100, 2, 100, 1, 2).
			WillReturnResult(sqlmock.NewResult(0, 2))

		mock.ExpectExec("INSERT INTO transactions").
			WithArgs(100, 2, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		err = repo.SendCoins(ctx, &domain.Transaction{
			Amount:     100,
			SenderID:   1,
			ReceiverID: 2,
		})

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUserRepo_BuyItem(t *testing.T) {
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := pgdb.NewUserRepo(db, pgdb.NewTransactionRepo(db), pgdb.NewInventoryRepo(db))

		mock.ExpectBegin()
		mock.ExpectQuery("SELECT price, id FROM merch").
			WithArgs("sword").
			WillReturnRows(sqlmock.NewRows([]string{"price", "id"}).AddRow(500, 1))

		mock.ExpectQuery("SELECT coins FROM users").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"coins"}).AddRow(1000))

		mock.ExpectExec("UPDATE users SET coins").
			WithArgs(500, 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectExec("INSERT INTO inventory").
			WithArgs(1, 1, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		err = repo.BuyItem(ctx, 1, "sword")
		assert.NoError(t, err)
	})

	t.Run("No records", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := pgdb.NewUserRepo(db, pgdb.NewTransactionRepo(db), pgdb.NewInventoryRepo(db))

		mock.ExpectBegin()
		mock.ExpectQuery("SELECT price, id FROM merch").
			WithArgs("invalid_item").
			WillReturnError(sql.ErrNoRows)

		mock.ExpectRollback()

		err = repo.BuyItem(ctx, 1, "invalid_item")
		assert.Error(t, err)
		assert.Equal(t, consts.ErrItemNotFound, err)
	})
}
