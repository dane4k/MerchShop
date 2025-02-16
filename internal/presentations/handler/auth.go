package handler

import (
	"errors"

	"github.com/dane4k/MerchShop/internal/application/service"
	"github.com/dane4k/MerchShop/internal/errs"
	"github.com/dane4k/MerchShop/internal/presentations/dto/request"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	userService service.UserService
	baseHandler *BaseHandler
}

func NewAuthHandler(userService service.UserService, baseHandler *BaseHandler) *AuthHandler {
	return &AuthHandler{userService: userService, baseHandler: baseHandler}
}

func (ah *AuthHandler) Auth(c *gin.Context) {
	var req request.AuthRequest
	if !ah.baseHandler.BindRequest(c, &req) {
		return
	}

	token, err := ah.userService.LoginUser(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		if errors.Is(err, errs.ErrInvalidPassword) {
			RespondWithError(c, 401, err.Error())
			return
		} else {
			RespondWithError(c, 400, err.Error())
			return
		}
	}

	c.JSON(200, gin.H{
		"token": token,
	})
}
