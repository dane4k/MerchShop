package middleware

import (
	"errors"
	"strings"

	"github.com/dane4k/MerchShop/internal/application/service"
	"github.com/dane4k/MerchShop/internal/presentations/handler"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	authService service.AuthService
}

func NewAuthMiddleware(authService service.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

func (amw *AuthMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := extractToken(c.GetHeader("Authorization"))
		if err != nil {
			handler.RespondWithError(c, 401, err.Error())
			return
		}

		userID, err := amw.authService.ParseJWT(token)
		if err != nil {
			handler.RespondWithError(c, 401, err.Error())
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}

func extractToken(header string) (string, error) {
	if header == "" {
		return "", errors.New("empty authorization header")
	}

	if !strings.HasPrefix(header, "Bearer ") {
		return "", errors.New("invalid authorization header")
	}

	return strings.TrimPrefix(header, "Bearer "), nil
}
