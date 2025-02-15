package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"strings"
)

type BaseHandler struct {
	validator *validator.Validate
}

func NewBaseHandler() *BaseHandler {
	return &BaseHandler{
		validator: validator.New(),
	}
}

func (bh *BaseHandler) BindRequest(c *gin.Context, req interface{}) bool {
	if err := c.ShouldBindJSON(req); err != nil {
		RespondWithError(c, 400, "invalid request body")
		return false
	}

	if err := bh.validator.Struct(req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errorMessages := make([]string, 0, len(validationErrors))

		for _, e := range validationErrors {
			errorMessages = append(errorMessages, fmt.Sprintf("%s: %s", e.Field(), e.Tag()))
		}

		errorMsg := "validation error: " + strings.Join(errorMessages, "; ")
		RespondWithError(c, 400, errorMsg)

		return false
	}

	return true
}

func RespondWithError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{
		"errors": message,
	})
}
