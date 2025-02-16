package service

import (
	"context"
	"errors"

	"github.com/dane4k/MerchShop/internal/domain"
	"github.com/dane4k/MerchShop/internal/errs"
	"github.com/dane4k/MerchShop/internal/infrastructure/db/pgdb"
	"github.com/dane4k/MerchShop/internal/presentations/dto/request"
	"github.com/dane4k/MerchShop/internal/presentations/dto/response"
)

type UserService interface {
	GetInfo(ctx context.Context, userID int) (*response.InfoResponse, error)
	SendCoins(ctx context.Context, fromUserID int, request *request.SendCoinRequest) error
	BuyItem(ctx context.Context, userID int, request request.BuyItemRequest) error
	LoginUser(ctx context.Context, username, password string) (string, error)
}

type UserRepo interface {
	GetUserData(ctx context.Context, username string) (*domain.User, error)
	AddUser(ctx context.Context, username string, hashedPassword string) (int, error)
	GetUserBalance(ctx context.Context, userID int) (int, error)
	SendCoins(ctx context.Context, transaction *domain.Transaction) error
	BuyItem(ctx context.Context, userID int, item string) error
}

type InventoryRepo interface {
	GetUserInventory(ctx context.Context, userID int) ([]*response.InventoryItem, error)
}

type CoinHistoryGetter interface {
	GetCoinHistory(ctx context.Context, userID int) (*response.CoinHistory, error)
}

type userService struct {
	userRepo           UserRepo
	inventoryRepo      InventoryRepo
	transactionService CoinHistoryGetter
	authService        AuthService
}

func NewUserService(userRepo UserRepo, inventoryRepo InventoryRepo,
	transactionService CoinHistoryGetter, authService AuthService) UserService {
	return &userService{
		userRepo:           userRepo,
		inventoryRepo:      inventoryRepo,
		transactionService: transactionService,
		authService:        authService,
	}
}

func (us *userService) LoginUser(ctx context.Context, username, password string) (string, error) {
	user, err := us.userRepo.GetUserData(ctx, username)
	if err != nil {
		if !(errors.Is(err, errs.ErrUserNotFound)) {
			return "", errs.ErrInternalServerError
		}
	}

	var userID int
	var hashedPassword string

	if user != nil {
		userID, hashedPassword = user.ID, user.PasswordHashed
	}

	if userID == 0 {
		hashedPassword, err = us.authService.EncryptPassword(password)
		if err != nil {
			return "", pgdb.RespondWithError(errs.ErrInternalServerError, "")
		}
		userID, err = us.userRepo.AddUser(ctx, username, hashedPassword)
		if err != nil {
			return "", pgdb.RespondWithError(errs.ErrInternalServerError, "")
		}
	} else {
		if err = us.authService.ComparePasswords(password, hashedPassword); err != nil {
			return "", errs.ErrInvalidPassword
		}
	}

	return us.authService.GenerateJWT(userID)
}

func (us *userService) GetInfo(ctx context.Context, userID int) (*response.InfoResponse, error) {
	coinsAmount, err := us.userRepo.GetUserBalance(ctx, userID)
	if err != nil {
		return nil, errs.ErrInternalServerError
	}

	userInventory, err := us.inventoryRepo.GetUserInventory(ctx, userID)
	if err != nil {
		return nil, errs.ErrInternalServerError
	}

	coinHistory, err := us.transactionService.GetCoinHistory(ctx, userID)
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return nil, errs.ErrUserNotFound
		}
		return nil, errs.ErrInternalServerError
	}

	return &response.InfoResponse{
		Coins:       coinsAmount,
		Inventory:   userInventory,
		CoinHistory: coinHistory,
	}, nil
}

func (us *userService) SendCoins(ctx context.Context, fromUserID int, request *request.SendCoinRequest) error {
	receiver, err := us.userRepo.GetUserData(ctx, request.ToUser)
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return errs.ErrWrongReceiver
		}
		return errs.ErrInternalServerError
	}

	err = us.userRepo.SendCoins(ctx, &domain.Transaction{
		Amount:     request.Amount,
		ReceiverID: receiver.ID,
		SenderID:   fromUserID,
	})
	if err != nil {
		if errors.Is(err, errs.ErrWrongReceiverID) {
			return errs.ErrWrongReceiver
		}
		if errors.Is(err, errs.ErrInsufficientFunds) {
			return err
		}
		return errs.ErrInternalServerError
	}

	return nil
}

func (us *userService) BuyItem(ctx context.Context, userID int, request request.BuyItemRequest) error {
	err := us.userRepo.BuyItem(ctx, userID, request.Name)
	if err != nil {
		if errors.Is(err, errs.ErrItemNotFound) {
			return err
		}
		if errors.Is(err, errs.ErrInsufficientFunds) {
			return err
		}
		return errs.ErrInternalServerError
	}

	return nil
}
