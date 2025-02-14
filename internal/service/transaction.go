package service

import (
	"MerchShop/internal/dto"
	"MerchShop/internal/repo/pgdb"
	"context"
	"errors"
)

type TransactionService interface {
	getCoinHistory(ctx context.Context, userID int) (*dto.CoinHistory, error)
}

type transactionService struct {
	userRepo        pgdb.UserRepo
	transactionRepo pgdb.TransactionRepo
}

func NewTransactionService(userRepo pgdb.UserRepo, transactionRepo pgdb.TransactionRepo) TransactionService {
	return &transactionService{userRepo: userRepo, transactionRepo: transactionRepo}
}

func (ts *transactionService) getCoinHistory(ctx context.Context, userID int) (*dto.CoinHistory, error) {
	username, err := ts.userRepo.GetNicknameByID(ctx, userID)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, ErrInternalServerError
	}

	transactions, err := ts.transactionRepo.GetUserTransactions(ctx, userID)
	if err != nil {
		return nil, ErrInternalServerError
	}

	if len(transactions) == 0 {
		return &dto.CoinHistory{
			Received: []*dto.ReceivedTransaction{},
			Sent:     []*dto.SentTransaction{},
		}, nil
	}

	var receivedTransactions []*dto.ReceivedTransaction
	var sentTransactions []*dto.SentTransaction

	for _, transaction := range transactions {
		if transaction.ReceiverUsername == username {
			receivedTransactions = append(receivedTransactions,
				&dto.ReceivedTransaction{
					FromUser: transaction.SenderUsername,
					Amount:   transaction.Amount,
				})
		} else {
			sentTransactions = append(sentTransactions,
				&dto.SentTransaction{
					ToUser: transaction.ReceiverUsername,
					Amount: transaction.Amount,
				})
		}
	}

	return &dto.CoinHistory{
		Received: receivedTransactions,
		Sent:     sentTransactions,
	}, nil
}
