package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dane4k/MerchShop/internal/application/service"
	"github.com/dane4k/MerchShop/internal/application/service_test/mocks"
	"github.com/dane4k/MerchShop/internal/domain"
	"github.com/dane4k/MerchShop/internal/presentations/dto/request"
	"github.com/dane4k/MerchShop/internal/presentations/handler"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthHandler_Auth(t *testing.T) {
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
		authHandler := handler.NewAuthHandler(userService, baseHandler)
		testUsername := "user1"
		testPassword := "12345678"
		hashedPassword := "hashed_pass123"

		userRepo.On("GetUserData", mock.Anything, testUsername).Return(&domain.User{
			ID:             1,
			PasswordHashed: hashedPassword,
		}, nil)
		authService.On("ComparePasswords", testPassword, hashedPassword).Return(nil)
		authService.On("GenerateJWT", 1).Return("token", nil)

		req := request.AuthRequest{
			Username: testUsername,
			Password: testPassword,
		}
		body, _ := json.Marshal(req)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/auth", bytes.NewBuffer(body))

		authHandler.Auth(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"token":"token"}`, w.Body.String())

		userRepo.AssertCalled(t, "GetUserData", mock.Anything, testUsername)
		authService.AssertCalled(t, "ComparePasswords", testPassword, hashedPassword)
		authService.AssertCalled(t, "GenerateJWT", 1)

		userRepo.AssertExpectations(t)
		authService.AssertExpectations(t)
	})
}
