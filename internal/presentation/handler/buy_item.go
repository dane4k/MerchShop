package handler

import (
	"errors"

	"github.com/dane4k/MerchShop/internal/application/service"
	"github.com/dane4k/MerchShop/internal/errs"
	"github.com/dane4k/MerchShop/internal/presentation/dto/request"
	"github.com/gin-gonic/gin"
)

type InventoryHandler struct {
	userService service.UserService
	baseHandler BaseHandler
}

func NewInventoryHandler(userService service.UserService, baseHandler *BaseHandler) *InventoryHandler {
	return &InventoryHandler{
		userService: userService,
		baseHandler: *baseHandler,
	}
}

func (ih *InventoryHandler) BuyItem(c *gin.Context) {
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

	item := c.Param("item")
	req := request.BuyItemRequest{Name: item}

	err := ih.userService.BuyItem(c.Request.Context(), uID, req)
	if err != nil {
		if errors.Is(err, errs.ErrItemNotFound) {
			RespondWithError(c, 404, err.Error())
			return
		} else if errors.Is(err, errs.ErrInsufficientFunds) {
			RespondWithError(c, 400, err.Error())
			return
		}
		RespondWithError(c, 500, err.Error())
		return
	}

	c.JSON(200, gin.H{})
}
