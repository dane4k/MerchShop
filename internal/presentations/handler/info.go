package handler

import (
	"errors"

	"github.com/dane4k/MerchShop/internal/application/service"
	"github.com/dane4k/MerchShop/internal/errs"
	"github.com/gin-gonic/gin"
)

type InfoHandler struct {
	userService service.UserService
}

func NewInfoHandler(userService service.UserService) *InfoHandler {
	return &InfoHandler{userService: userService}
}

func (ih *InfoHandler) GetInfo(c *gin.Context) {
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

	info, err := ih.userService.GetInfo(c.Request.Context(), uID)
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			RespondWithError(c, 400, err.Error())
			return
		}
		RespondWithError(c, 500, err.Error())
		return
	}

	c.JSON(200, info)
}
