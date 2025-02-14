package service

import (
	"MerchShop/internal/dto/request"
	"MerchShop/internal/dto/response"
	"MerchShop/internal/entity"
	"MerchShop/internal/repo/pgdb"
	"context"
	"errors"
)

type UserService interface {
	GetInfo(ctx context.Context, userID int) (*response.InfoResponse, error)
	SendCoins(ctx context.Context, fromUserID int, request *request.SendCoinRequest) error
	BuyItem(ctx context.Context, userID int, request request.BuyItemRequest) error
	LoginUser(ctx context.Context, username, password string) (string, error)
}

type userService struct {
	userRepo           pgdb.UserRepo
	inventoryRepo      pgdb.InventoryRepo
	transactionService TransactionService
	authService        AuthService
}

func NewUserService(userRepo pgdb.UserRepo, inventoryRepo pgdb.InventoryRepo,
	transactionService TransactionService, authService AuthService) UserService {
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
		if !(errors.Is(err, pgdb.ErrUserNotFound)) {
			return "", pgdb.ErrInternalServerError
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
			return "", pgdb.RespondWithError(pgdb.ErrInternalServerError, "")
		}
		userID, err = us.userRepo.AddUser(ctx, username, hashedPassword)
		if err != nil {
			return "", pgdb.RespondWithError(pgdb.ErrInternalServerError, "")
		}
	} else {
		if err = us.authService.ComparePasswords(password, hashedPassword); err != nil {
			return "", ErrInvalidPassword
		}
	}

	return us.authService.GenerateJWT(userID)
}

func (us *userService) GetInfo(ctx context.Context, userID int) (*response.InfoResponse, error) {
	coinsAmount, err := us.userRepo.GetUserBalance(ctx, userID)
	if err != nil {
		return nil, pgdb.ErrInternalServerError
	}

	userInventory, err := us.inventoryRepo.GetUserInventory(ctx, userID)
	if err != nil {
		return nil, pgdb.ErrInternalServerError
	}

	coinHistory, err := us.transactionService.getCoinHistory(ctx, userID)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, pgdb.ErrUserNotFound
		}
		return nil, pgdb.ErrInternalServerError
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
		if errors.Is(err, pgdb.ErrUserNotFound) {
			return ErrWrongReceiver
		}
		return ErrInternalServerError
	}

	err = us.userRepo.SendCoins(ctx, &entity.Transaction{
		Amount:     request.Amount,
		ReceiverID: receiver.ID,
		SenderID:   fromUserID,
	})
	if err != nil {
		if errors.Is(err, pgdb.ErrWrongReceiverID) {
			return ErrWrongReceiver
		}
		if errors.Is(err, pgdb.ErrInsufficientFunds) {
			return err
		}
		return ErrInternalServerError
	}

	return nil
}

func (us *userService) BuyItem(ctx context.Context, userID int, request request.BuyItemRequest) error {
	err := us.userRepo.BuyItem(ctx, userID, request.Name)
	if err != nil {
		if errors.Is(err, pgdb.ErrItemNotFound) {
			return err
		}
		if errors.Is(err, pgdb.ErrInsufficientFunds) {
			return err
		}
		return ErrInternalServerError
	}

	return nil
}
