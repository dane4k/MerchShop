package service

import (
	"context"
	"errors"

	"github.com/dane4k/MerchShop/internal/errs"
	"github.com/dane4k/MerchShop/internal/presentations/dto/response"
)

type TransactionService interface {
	GetCoinHistory(ctx context.Context, userID int) (*response.CoinHistory, error)
}

type TransactionRepo interface {
	GetUserTransactions(ctx context.Context, userID int) ([]*response.Transaction, error)
}

type UserRepoForTransactions interface {
	GetNicknameByID(ctx context.Context, userID int) (string, error)
}

type transactionService struct {
	userRepo        UserRepoForTransactions
	transactionRepo TransactionRepo
}

func NewTransactionService(userRepo UserRepoForTransactions, transactionRepo TransactionRepo) *transactionService {
	return &transactionService{userRepo: userRepo, transactionRepo: transactionRepo}
}

func (ts *transactionService) GetCoinHistory(ctx context.Context, userID int) (*response.CoinHistory, error) {
	username, err := ts.userRepo.GetNicknameByID(ctx, userID)
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return nil, errs.ErrUserNotFound
		}
		return nil, errs.ErrInternalServerError
	}

	transactions, err := ts.transactionRepo.GetUserTransactions(ctx, userID)
	if err != nil {
		return nil, errs.ErrInternalServerError
	}

	if len(transactions) == 0 {
		return &response.CoinHistory{
			Received: []*response.ReceivedTransaction{},
			Sent:     []*response.SentTransaction{},
		}, nil
	}

	var receivedTransactions []*response.ReceivedTransaction
	var sentTransactions []*response.SentTransaction

	for _, transaction := range transactions {
		if transaction.ReceiverUsername == username {
			receivedTransactions = append(receivedTransactions,
				&response.ReceivedTransaction{
					FromUser: transaction.SenderUsername,
					Amount:   transaction.Amount,
				})
		} else {
			sentTransactions = append(sentTransactions,
				&response.SentTransaction{
					ToUser: transaction.ReceiverUsername,
					Amount: transaction.Amount,
				})
		}
	}

	return &response.CoinHistory{
		Received: receivedTransactions,
		Sent:     sentTransactions,
	}, nil
}
