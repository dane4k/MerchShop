package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/dane4k/MerchShop/domain"
)

type UserRepo struct {
	mock.Mock
}

func (urm *UserRepo) GetUserData(ctx context.Context, username string) (*domain.User, error) {
	args := urm.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (urm *UserRepo) AddUser(ctx context.Context, username string, hashedPassword string) (int, error) {
	args := urm.Called(ctx, username, hashedPassword)
	return args.Int(0), args.Error(1)
}

func (urm *UserRepo) GetUserBalance(ctx context.Context, userID int) (int, error) {
	args := urm.Called(ctx, userID)
	return args.Int(0), args.Error(1)
}

func (urm *UserRepo) SendCoins(ctx context.Context, transaction *domain.Transaction) error {
	args := urm.Called(ctx, transaction)
	return args.Error(0)
}

func (urm *UserRepo) BuyItem(ctx context.Context, userID int, item string) error {
	args := urm.Called(ctx, userID, item)
	return args.Error(0)
}
