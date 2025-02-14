package service

import "errors"

var (
	ErrUserNotFound        = errors.New("user not found")
	ErrInvalidPassword     = errors.New("invalid password")
	ErrInternalServerError = errors.New("internal server error")
	ErrWrongReceiver       = errors.New("wrong receiver")
)
