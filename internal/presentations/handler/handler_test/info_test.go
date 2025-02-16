package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dane4k/MerchShop/internal/application/service"
	"github.com/dane4k/MerchShop/internal/application/service_test/mocks"
	"github.com/dane4k/MerchShop/internal/presentations/dto/response"
	"github.com/dane4k/MerchShop/internal/presentations/handler"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestInfoHandler_GetInfo(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		userRepo := new(mocks.UserRepo)
		inventoryRepo := new(mocks.InventoryRepo)
		transactionService := new(mocks.TransactionService)
		authService := new(mocks.AuthService)

		userService := service.NewUserService(
			userRepo,
			inventoryRepo,
			transactionService,
			authService,
		)
		infoHandler := handler.NewInfoHandler(userService)

		userRepo.On("GetUserBalance", mock.Anything, 1).
			Return(1000, nil)
		inventoryRepo.On("GetUserInventory", mock.Anything, 1).
			Return([]*response.InventoryItem{
				{Type: "umbrella", Quantity: 1}}, nil)
		transactionService.On("GetCoinHistory", mock.Anything, 1).
			Return(&response.CoinHistory{}, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/info", nil)
		c.Set("userID", 1)

		infoHandler.GetInfo(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"coins":1000`)
		userRepo.AssertExpectations(t)
		inventoryRepo.AssertExpectations(t)
		transactionService.AssertExpectations(t)
	})
}
