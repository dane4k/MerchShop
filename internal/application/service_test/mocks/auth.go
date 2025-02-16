package mocks

import "github.com/stretchr/testify/mock"

type AuthService struct {
	mock.Mock
}

func (asm *AuthService) GenerateJWT(userID int) (string, error) {
	args := asm.Called(userID)
	return args.String(0), args.Error(1)
}

func (asm *AuthService) ParseJWT(signedJWT string) (int, error) {
	args := asm.Called(signedJWT)
	return args.Int(0), args.Error(1)
}

func (asm *AuthService) EncryptPassword(password string) (string, error) {
	args := asm.Called(password)
	return args.String(0), args.Error(1)
}

func (asm *AuthService) ComparePasswords(password string, hashedPassword string) error {
	args := asm.Called(password, hashedPassword)
	return args.Error(0)
}
