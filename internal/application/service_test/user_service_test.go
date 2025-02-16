package service_test

import (
	"context"
	"testing"

	"github.com/dane4k/MerchShop/internal/application/service"
	"github.com/dane4k/MerchShop/internal/application/service_test/mocks"
	"github.com/dane4k/MerchShop/internal/domain"
	"github.com/dane4k/MerchShop/internal/errs"
	"github.com/dane4k/MerchShop/internal/presentation/dto/request"
	"github.com/dane4k/MerchShop/internal/presentation/dto/response"
	"github.com/stretchr/testify/assert"
)

func TestUserService(t *testing.T) {
	ctx := context.Background()

	t.Run("Successful user login", func(t *testing.T) {
		userRepo := new(mocks.UserRepo)
		auth := new(mocks.AuthService)
		us := service.NewUserService(userRepo, nil, nil, auth)

		userRepo.On("GetUserData", ctx, "user").Return(&domain.User{
			ID:             1,
			PasswordHashed: "hash",
		}, nil)
		auth.On("ComparePasswords", "pass", "hash").
			Return(nil)
		auth.On("GenerateJWT", 1).
			Return("token", nil)

		token, err := us.LoginUser(ctx, "user", "pass")

		assert.NoError(t, err)
		assert.Equal(t, "token", token)
		userRepo.AssertExpectations(t)
	})

	t.Run("Successful new user login", func(t *testing.T) {
		userRepo := new(mocks.UserRepo)
		auth := new(mocks.AuthService)
		us := service.NewUserService(userRepo, nil, nil, auth)

		userRepo.On("GetUserData", ctx, "new_user").
			Return(nil, errs.ErrUserNotFound)
		auth.On("EncryptPassword", "pass").
			Return("encrypted", nil)
		userRepo.On("AddUser", ctx, "new_user", "encrypted").
			Return(2, nil)
		auth.On("GenerateJWT", 2).
			Return("new_token", nil)

		token, err := us.LoginUser(ctx, "new_user", "pass")

		assert.NoError(t, err)
		assert.Equal(t, "new_token", token)
		userRepo.AssertExpectations(t)
	})

	t.Run("GetInfoSuccess", func(t *testing.T) {
		userRepo := new(mocks.UserRepo)
		inventory := new(mocks.InventoryRepo)
		tx := new(mocks.TransactionService)
		us := service.NewUserService(userRepo, inventory, tx, nil)

		userRepo.On("GetUserBalance", ctx, 1).
			Return(1000, nil)
		inventory.On("GetUserInventory", ctx, 1).
			Return([]*response.InventoryItem{
				{Type: "sword", Quantity: 1},
			}, nil)
		tx.On("GetCoinHistory", ctx, 1).
			Return(&response.CoinHistory{}, nil)

		info, err := us.GetInfo(ctx, 1)

		assert.NoError(t, err)
		assert.Equal(t, 1000, info.Coins)
		assert.Len(t, info.Inventory, 1)
	})

	t.Run("SendCoins Error", func(t *testing.T) {
		userRepo := new(mocks.UserRepo)
		us := service.NewUserService(userRepo, nil, nil, nil)

		userRepo.On("GetUserData", ctx, "unknown").
			Return(nil, errs.ErrUserNotFound)

		err := us.SendCoins(ctx, 1, &request.SendCoinRequest{
			ToUser: "unknown",
			Amount: 100,
		})

		assert.Error(t, err)
		userRepo.AssertExpectations(t)
	})

	t.Run("BuyItem Error", func(t *testing.T) {
		userRepo := new(mocks.UserRepo)
		us := service.NewUserService(userRepo, nil, nil, nil)

		userRepo.On("BuyItem", ctx, 1, "item").
			Return(errs.ErrItemNotFound)

		err := us.BuyItem(ctx, 1, request.BuyItemRequest{Name: "item"})

		assert.Error(t, err)
		userRepo.AssertExpectations(t)
	})
}
