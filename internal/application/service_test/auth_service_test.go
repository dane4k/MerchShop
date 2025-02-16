package service_test

import (
	"testing"
	"time"

	"github.com/dane4k/MerchShop/internal/application/service"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestAuthService_GenerateParseJWT(t *testing.T) {
	as := service.NewAuthService("secret")
	userID := 123

	t.Run("Success", func(t *testing.T) {
		token, err := as.GenerateJWT(userID)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		parsedID, err := as.ParseJWT(token)
		assert.NoError(t, err)
		assert.Equal(t, userID, parsedID)
	})

	t.Run("Incorrect (expired) JWT token", func(t *testing.T) {
		expiredJWT := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
			"user_id": userID,
			"exp":     time.Now().Add(-time.Second).Unix(),
		})

		expiredToken, err := expiredJWT.SignedString([]byte("secret"))
		if err != nil {
			t.Fatal(err)
		}

		_, err = as.ParseJWT(expiredToken)
		assert.Error(t, err)
	})
}
