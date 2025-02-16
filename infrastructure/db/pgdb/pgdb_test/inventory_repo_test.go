package pgdb_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/dane4k/MerchShop/infrastructure/db/pgdb"
)

func TestInventoryRepo_AddItem(t *testing.T) {
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := pgdb.NewInventoryRepo(db)

		mock.ExpectBegin()
		tx, err := db.Begin()
		assert.NoError(t, err)

		mock.ExpectExec("INSERT INTO inventory").WithArgs(1, 1, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err = repo.AddItem(ctx, tx, 1, 1)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
func TestInventoryRepo_GetUserInventory(t *testing.T) {
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := pgdb.NewInventoryRepo(db)
		rows := sqlmock.NewRows([]string{"name", "quantity"}).
			AddRow("sword", 2).
			AddRow("shield", 1)

		mock.ExpectQuery("SELECT m.name, quantity FROM inventory i JOIN merch m").
			WithArgs(1).
			WillReturnRows(rows)

		items, err := repo.GetUserInventory(ctx, 1)
		assert.NoError(t, err)
		assert.Len(t, items, 2)
	})

	t.Run("No records", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := pgdb.NewInventoryRepo(db)
		mock.ExpectQuery("SELECT m.name, quantity FROM inventory i JOIN merch m").
			WithArgs(1).
			WillReturnError(sql.ErrNoRows)

		_, err = repo.GetUserInventory(ctx, 1)
		assert.Error(t, err)
	})
}
