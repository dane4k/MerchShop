package pgdb

import (
	"context"
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"

	"github.com/dane4k/MerchShop/internal/consts"
	"github.com/dane4k/MerchShop/presentations/dto/response"
)

type InventoryRepo struct {
	SQL sq.StatementBuilderType
}

func NewInventoryRepo(DB *sql.DB) *InventoryRepo {
	return &InventoryRepo{
		SQL: sq.StatementBuilder.PlaceholderFormat(sq.Dollar).RunWith(DB),
	}
}

func (ir *InventoryRepo) AddItem(ctx context.Context, tx *sql.Tx, userID, merchID int) error {
	if err := checkTx(tx); err != nil {
		return RespondWithError(err, "error adding item")
	}

	res, err := ir.SQL.Insert(consts.InventoryTable).
		Columns(consts.ColumnUserID, consts.ColumnMerchID, consts.ColumnQuantity).
		Values(userID, merchID, 1).
		Suffix("ON CONFLICT (user_id, merch_id) DO UPDATE SET quantity = inventory.quantity + 1").
		RunWith(tx).
		ExecContext(ctx)
	if err != nil {
		return RespondWithError(err, "error adding item")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return RespondWithError(err, "error getting rows affected")
	}
	if rowsAffected == 0 {
		return RespondWithError(errors.New("unable to add item"), "error adding item")
	}

	return nil
}

func (ir *InventoryRepo) GetUserInventory(ctx context.Context, userID int) ([]*response.InventoryItem, error) {
	rows, err := ir.SQL.Select("m.name", consts.ColumnQuantity).
		From("inventory i").
		Join("merch m on i.merch_id = m.id").
		Where(sq.Eq{"i.user_id": userID}).
		QueryContext(ctx)
	if err != nil {
		return nil, RespondWithError(err, "error getting user inventory")
	}
	defer rows.Close()

	var items []*response.InventoryItem
	for rows.Next() {
		var item response.InventoryItem
		err = rows.Scan(&item.Type, &item.Quantity)
		if err != nil {
			return nil, RespondWithError(err, "error scanning inventory row")
		}
		items = append(items, &item)
	}

	if err = rows.Err(); err != nil {
		return nil, RespondWithError(err, "inventory rows error")
	}

	return items, nil
}

func checkTx(tx *sql.Tx) error {
	if tx == nil {
		return errors.New("tx is nil")
	}
	return nil
}
