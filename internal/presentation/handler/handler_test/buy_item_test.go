package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dane4k/MerchShop/internal/application/service"
	"github.com/dane4k/MerchShop/internal/application/service_test/mocks"
	"github.com/dane4k/MerchShop/internal/presentation/handler"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBuyHandler_BuyItem(t *testing.T) {
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
		inventoryHandler := handler.NewInventoryHandler(userService, baseHandler)

		userRepo.On("BuyItem", mock.Anything, 1, "umbrella").
			Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/buy/umbrella", nil)
		c.Params = gin.Params{{Key: "item", Value: "umbrella"}}
		c.Set("userID", 1)

		inventoryHandler.BuyItem(c)

		assert.Equal(t, http.StatusOK, w.Code)
		userRepo.AssertExpectations(t)
	})
}
