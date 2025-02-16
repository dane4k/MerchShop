package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/dane4k/MerchShop/presentations/dto/response"
)

type TransactionService struct {
	mock.Mock
}

type UserRepoForTransactions struct {
	mock.Mock
}

type TransactionRepo struct {
	mock.Mock
}

func (tsm *TransactionService) GetCoinHistory(ctx context.Context, userID int) (*response.CoinHistory, error) {
	args := tsm.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.CoinHistory), args.Error(1)
}

func (utm *UserRepoForTransactions) GetNicknameByID(ctx context.Context, userID int) (string, error) {
	args := utm.Called(ctx, userID)
	return args.String(0), args.Error(1)
}

func (trm *TransactionRepo) GetUserTransactions(ctx context.Context, userID int) ([]*response.Transaction, error) {
	args := trm.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*response.Transaction), args.Error(1)
}
