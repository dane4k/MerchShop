package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/dane4k/MerchShop/infrastructure/db/pgdb"
)

type AuthService interface {
	GenerateJWT(userID int) (string, error)
	ParseJWT(signedJWT string) (int, error)
	EncryptPassword(password string) (string, error)
	ComparePasswords(password string, hashedPassword string) error
}

const (
	TokenLT = 24 * time.Hour
)

type authService struct {
	JWTSecret string
}

func NewAuthService(secret string) AuthService {
	return &authService{JWTSecret: secret}
}

func (as *authService) GenerateJWT(userID int) (string, error) {
	jwtToken := jwt.NewWithClaims(
		jwt.SigningMethodHS512,
		jwt.MapClaims{
			"user_id": userID,
			"exp":     time.Now().Add(TokenLT).Unix(),
		})

	token, err := jwtToken.SignedString([]byte(as.JWTSecret))
	if err != nil {
		return "", pgdb.RespondWithError(errors.New("internal server error"), "error generating JWT")
	}

	return token, nil
}

func (as *authService) ParseJWT(signedJWT string) (int, error) {
	token, err := jwt.Parse(signedJWT, func(token *jwt.Token) (interface{}, error) {
		return []byte(as.JWTSecret), nil
	})
	if err != nil || !token.Valid {
		return 0, pgdb.RespondWithError(errors.New("invalid token"), "error parsing JWT")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, pgdb.RespondWithError(errors.New("invalid token claims"), "error parsing JWT")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, pgdb.RespondWithError(errors.New("invalid user ID in token"), "error parsing JWT")
	}

	return int(userID), nil
}

func (as *authService) EncryptPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", pgdb.RespondWithError(err, "error encrypting password")
	}
	return string(hashed), err
}

func (as *authService) ComparePasswords(password string, hashedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return pgdb.RespondWithError(err, "error comparing passwords")
	}
	return nil
}
