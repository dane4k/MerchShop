package handler

import (
	"MerchShop/internal/dto/request"
	"MerchShop/internal/service"
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
		RespondWithError(c, 400, err.Error())
		return
	}

	c.JSON(200, gin.H{
		"token": token,
	})
}
