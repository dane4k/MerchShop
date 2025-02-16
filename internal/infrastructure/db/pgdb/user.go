package pgdb

import (
	"context"
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/dane4k/MerchShop/internal/domain"
	"github.com/dane4k/MerchShop/internal/errs"
	"github.com/sirupsen/logrus"
)

type UserRepo struct {
	DB              *sql.DB
	SQL             sq.StatementBuilderType
	transactionRepo *TransactionRepo
	inventoryRepo   *InventoryRepo
}

func NewUserRepo(DB *sql.DB, transactionRepo *TransactionRepo, inventoryRepo *InventoryRepo) *UserRepo {
	return &UserRepo{
		DB:              DB,
		SQL:             sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		transactionRepo: transactionRepo,
		inventoryRepo:   inventoryRepo,
	}
}

func (ur *UserRepo) GetNicknameByID(ctx context.Context, userID int) (string, error) {
	var username string
	err := ur.SQL.Select(ColumnUsername).
		From(UsersTable).
		Where(sq.Eq{"id": userID}).
		RunWith(ur.DB).
		QueryRowContext(ctx).
		Scan(&username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", RespondWithError(errs.ErrUserNotFound, "user not found")
		}
		return "", RespondWithError(err, "failed to get user")
	}
	return username, nil
}

func (ur *UserRepo) GetUserData(ctx context.Context, username string) (*domain.User, error) {
	var userID int
	var hashedPassword string

	err := ur.SQL.Select("id", ColumnPasswordHashed).
		From(UsersTable).
		Where(sq.Eq{ColumnUsername: username}).
		RunWith(ur.DB).
		QueryRowContext(ctx).
		Scan(&userID, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, RespondWithError(errs.ErrUserNotFound, "error getting user data")
		}
		return nil, RespondWithError(err, "error getting user data")
	}

	return &domain.User{
		ID:             userID,
		PasswordHashed: hashedPassword,
	}, nil
}

func (ur *UserRepo) AddUser(ctx context.Context, username string, hashedPassword string) (int, error) {
	var userID int

	err := ur.SQL.Insert(UsersTable).
		Columns(ColumnUsername, ColumnCoins, ColumnPasswordHashed).
		Values(username, 1000, hashedPassword).
		Suffix("ON CONFLICT (username) DO NOTHING").
		Suffix("RETURNING id").
		RunWith(ur.DB).
		QueryRowContext(ctx).
		Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			user, err := ur.GetUserData(ctx, username)
			if err != nil {
				return 0, err
			}
			return user.ID, nil
		}
		return 0, RespondWithError(err, "error adding user")
	}

	return userID, nil
}

func (ur *UserRepo) GetUserBalance(ctx context.Context, userID int) (int, error) {
	return ur.getUserBalance(ctx, nil, userID)
}

func (ur *UserRepo) SendCoins(ctx context.Context, transaction *domain.Transaction) error {
	tx, err := ur.DB.Begin()
	if err != nil {
		return RespondWithError(err, "failed to start transaction")
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				logrus.WithError(rollbackErr).Error("failed to rollback transaction")
			}
		}
	}()

	balance, err := ur.getUserBalance(ctx, tx, transaction.SenderID)
	if err != nil {
		return RespondWithError(err, "failed to get sender balance")
	}
	if balance < transaction.Amount {
		return errs.ErrInsufficientFunds
	}

	res, err := ur.SQL.Update(UsersTable).
		Set(ColumnCoins, sq.Case().
			When(sq.Eq{"id": transaction.SenderID}, sq.Expr("coins - ?", transaction.Amount)).
			When(sq.Eq{"id": transaction.ReceiverID}, sq.Expr("coins + ?", transaction.Amount))).
		Where(sq.Or{
			sq.Eq{"id": transaction.SenderID},
			sq.Eq{"id": transaction.ReceiverID},
		}).
		RunWith(tx).
		ExecContext(ctx)
	if err != nil {
		return RespondWithError(err, "error updating user balance")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return RespondWithError(err, "error getting rows affected")
	}
	if rowsAffected != 2 {
		return RespondWithError(errs.ErrWrongReceiverID, "error sending coins")
	}

	err = ur.transactionRepo.AddTransaction(ctx, tx, transaction)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return RespondWithError(err, "error committing transaction")
	}

	return nil
}

func (ur *UserRepo) BuyItem(ctx context.Context, userID int, item string) error {
	tx, err := ur.DB.Begin()
	if err != nil {
		return RespondWithError(err, "failed to start transaction")
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				logrus.WithError(rollbackErr).Error("failed to rollback transaction")
			}
		}
	}()

	var price, itemID int
	err = ur.SQL.Select(ColumnPrice, "id").
		From(MerchTable).
		Where(sq.Eq{ColumnName: item}).
		RunWith(tx).
		QueryRowContext(ctx).
		Scan(&price, &itemID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return RespondWithError(errs.ErrItemNotFound, "error buying item")
		}
		return RespondWithError(err, "error buying item")
	}

	balance, err := ur.getUserBalance(ctx, tx, userID)
	if err != nil {
		return err
	}
	if balance < price {
		return errs.ErrInsufficientFunds
	}

	res, err := ur.SQL.Update(UsersTable).
		Set(ColumnCoins, sq.Expr("coins - ?", price)).
		Where(sq.Eq{"id": userID}).
		RunWith(tx).
		ExecContext(ctx)
	if err != nil {
		return RespondWithError(err, "error updating user balance")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return RespondWithError(err, "error getting rows affected")
	}
	if rowsAffected == 0 {
		return RespondWithError(errs.ErrUnableToUpdate, "error buying item")
	}

	err = ur.inventoryRepo.AddItem(ctx, tx, userID, itemID)
	if err != nil {
		return RespondWithError(err, "error adding item")
	}

	if err = tx.Commit(); err != nil {
		return RespondWithError(err, "error committing transaction")
	}

	return nil
}

func (ur *UserRepo) getUserBalance(ctx context.Context, tx *sql.Tx, userID int) (int, error) {
	query := ur.SQL.Select(ColumnCoins).
		From(UsersTable).
		Where(sq.Eq{"id": userID})

	var coins int
	if tx != nil {
		err := query.
			//Suffix("FOR UPDATE").
			RunWith(tx).
			QueryRowContext(ctx).
			Scan(&coins)
		return coins, err
	} else {
		err := query.RunWith(ur.DB).
			QueryRowContext(ctx).
			Scan(&coins)
		return coins, err
	}
}

func RespondWithError(err error, message string) error {
	logrus.WithError(err).Error(message)
	return err
}
