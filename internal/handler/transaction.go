package handler

import (
	"MerchShop/internal/dto/request"
	"MerchShop/internal/repo/pgdb"
	"MerchShop/internal/service"
	"errors"
	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	userService service.UserService
	baseHandler BaseHandler
}

func NewTransactionHandler(userService service.UserService, baseHandler *BaseHandler) *TransactionHandler {
	return &TransactionHandler{
		userService: userService,
		baseHandler: *baseHandler,
	}
}

func (th *TransactionHandler) SendCoins(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok {
		RespondWithError(c, 401, "unauthorized")
		return
	}

	uID, ok := userID.(int)
	if !ok {
		RespondWithError(c, 500, "internal server error")
		return
	}

	var req request.SendCoinRequest
	if !th.baseHandler.BindRequest(c, &req) {
		return
	}

	err := th.userService.SendCoins(c.Request.Context(), uID, &req)
	if err != nil {
		if errors.Is(err, service.ErrWrongReceiver) || errors.Is(err, pgdb.ErrInsufficientFunds) {
			RespondWithError(c, 400, err.Error())
			return
		}
		RespondWithError(c, 500, err.Error())
		return
	}

	c.JSON(200, gin.H{})
}
