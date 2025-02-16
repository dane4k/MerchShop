package service_test

import (
	"context"
	"testing"

	"github.com/dane4k/MerchShop/internal/application/service"
	"github.com/dane4k/MerchShop/internal/application/service_test/mocks"
	"github.com/dane4k/MerchShop/internal/errs"
	"github.com/dane4k/MerchShop/internal/presentations/dto/response"
	"github.com/stretchr/testify/assert"
)

func TestTransactionService(t *testing.T) {
	ctx := context.Background()
	userID := 1

	t.Run("Success", func(t *testing.T) {
		mockUserRepo := new(mocks.UserRepoForTransactions)
		mockTransactionRepo := new(mocks.TransactionRepo)
		ts := service.NewTransactionService(mockUserRepo, mockTransactionRepo)

		mockUserRepo.On("GetNicknameByID", ctx, userID).
			Return("user1", nil)
		mockTransactionRepo.On("GetUserTransactions", ctx, userID).
			Return([]*response.Transaction{
				{SenderUsername: "user", ReceiverUsername: "user1", Amount: 111},
				{SenderUsername: "user1", ReceiverUsername: "user3", Amount: 11},
			}, nil)

		result, err := ts.GetCoinHistory(ctx, userID)

		assert.NoError(t, err)
		assert.Equal(t, 1, len(result.Received))
		assert.Equal(t, 1, len(result.Sent))
		mockUserRepo.AssertExpectations(t)
		mockTransactionRepo.AssertExpectations(t)
	})

	t.Run("User doesnt exist", func(t *testing.T) {
		mockUserRepo := new(mocks.UserRepoForTransactions)
		mockTransactionRepo := new(mocks.TransactionRepo)
		ts := service.NewTransactionService(mockUserRepo, mockTransactionRepo)

		mockUserRepo.On("GetNicknameByID", ctx, 999).
			Return("", errs.ErrUserNotFound)

		_, err := ts.GetCoinHistory(ctx, 999)

		assert.Error(t, err)
		assert.ErrorIs(t, err, errs.ErrUserNotFound)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("No transactions", func(t *testing.T) {
		mockUserRepo := new(mocks.UserRepoForTransactions)
		mockTransactionRepo := new(mocks.TransactionRepo)
		ts := service.NewTransactionService(mockUserRepo, mockTransactionRepo)

		mockUserRepo.On("GetNicknameByID", ctx, userID).
			Return("user1", nil)
		mockTransactionRepo.On("GetUserTransactions", ctx, userID).
			Return([]*response.Transaction{}, nil)

		result, err := ts.GetCoinHistory(ctx, userID)

		assert.NoError(t, err)
		assert.Empty(t, result.Received)
		assert.Empty(t, result.Sent)
		mockUserRepo.AssertExpectations(t)
		mockTransactionRepo.AssertExpectations(t)
	})

	t.Run("Internal server error", func(t *testing.T) {
		mockUserRepo := new(mocks.UserRepoForTransactions)
		mockTransactionRepo := new(mocks.TransactionRepo)
		ts := service.NewTransactionService(mockUserRepo, mockTransactionRepo)

		mockUserRepo.On("GetNicknameByID", ctx, userID).
			Return("user1", nil)
		mockTransactionRepo.On("GetUserTransactions", ctx, userID).
			Return(nil, errs.ErrInternalServerError)

		_, err := ts.GetCoinHistory(ctx, userID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, errs.ErrInternalServerError)
		mockUserRepo.AssertExpectations(t)
		mockTransactionRepo.AssertExpectations(t)
	})
}
