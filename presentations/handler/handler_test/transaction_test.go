package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/dane4k/MerchShop/application/service"
	"github.com/dane4k/MerchShop/application/service_test/mocks"
	"github.com/dane4k/MerchShop/domain"
	"github.com/dane4k/MerchShop/internal/consts"
	"github.com/dane4k/MerchShop/presentations/dto/request"
	"github.com/dane4k/MerchShop/presentations/handler"
)

func TestTransactionHandler_SendCoins(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		userRepo := new(mocks.UserRepo)
		authService := new(mocks.AuthService)
		userService := service.NewUserService(
			userRepo,
			nil,
			nil,
			authService,
		)
		baseHandler := handler.NewBaseHandler()
		transactionHandler := handler.NewTransactionHandler(userService, baseHandler)
		userRepo.On("GetUserData", mock.Anything, "receiver").Return(
			&domain.User{ID: 2},
			nil,
		)

		userRepo.On("SendCoins", mock.Anything, mock.MatchedBy(func(tx *domain.Transaction) bool {
			return tx.SenderID == 1 && tx.ReceiverID == 2 && tx.Amount == 100
		})).
			Return(nil)

		req := request.SendCoinRequest{
			ToUser: "receiver",
			Amount: 100,
		}
		body, _ := json.Marshal(req)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/send", bytes.NewBuffer(body))
		c.Set("userID", 1)

		transactionHandler.SendCoins(c)

		assert.Equal(t, http.StatusOK, w.Code)
		userRepo.AssertExpectations(t)
	})

	t.Run("Not success: receiver not found", func(t *testing.T) {
		userRepo := new(mocks.UserRepo)
		authService := new(mocks.AuthService)
		userService := service.NewUserService(
			userRepo,
			nil,
			nil,
			authService,
		)
		baseHandler := handler.NewBaseHandler()
		transactionHandler := handler.NewTransactionHandler(userService, baseHandler)

		userRepo.On("GetUserData", mock.Anything, "user9999").Return(
			nil,
			consts.ErrUserNotFound,
		)

		req := request.SendCoinRequest{
			ToUser: "user9999",
			Amount: 100,
		}
		body, _ := json.Marshal(req)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/send", bytes.NewBuffer(body))
		c.Set("userID", 1)

		transactionHandler.SendCoins(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		userRepo.AssertExpectations(t)
	})

	t.Run("Insufficient funds", func(t *testing.T) {
		userRepo := new(mocks.UserRepo)
		authService := new(mocks.AuthService)
		userService := service.NewUserService(
			userRepo,
			nil,
			nil,
			authService,
		)
		baseHandler := handler.NewBaseHandler()
		transactionHandler := handler.NewTransactionHandler(userService, baseHandler)

		userRepo.On("GetUserData", mock.Anything, "receiver").
			Return(&domain.User{ID: 2}, nil)
		userRepo.On("SendCoins", mock.Anything, mock.Anything).Return(consts.ErrInsufficientFunds)

		req := request.SendCoinRequest{
			ToUser: "receiver",
			Amount: 1000,
		}
		body, _ := json.Marshal(req)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/send", bytes.NewBuffer(body))
		c.Set("userID", 1)

		transactionHandler.SendCoins(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "insufficient funds")
		userRepo.AssertExpectations(t)
	})
}
