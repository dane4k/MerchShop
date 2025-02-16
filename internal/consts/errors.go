package consts

import "errors"

var (
	ErrInsufficientFunds   = errors.New("insufficient funds")
	ErrUserNotFound        = errors.New("user not found")
	ErrItemNotFound        = errors.New("item not found")
	ErrWrongReceiverID     = errors.New("wrong receiver id")
	ErrUnableToUpdate      = errors.New("unable to update balance")
	ErrInternalServerError = errors.New("internal server error")
	ErrInvalidPassword     = errors.New("invalid password")
	ErrWrongReceiver       = errors.New("wrong receiver")
)
