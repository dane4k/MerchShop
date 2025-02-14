package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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
		errors := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			field, tag := err.Field(), err.Tag()
			errors[field] = tag
		}
		RespondWithError(c, 400, "validation failed", gin.H{"details": errors})
		return false
	}

	return true
}

func RespondWithError(c *gin.Context, statusCode int, message string, details ...gin.H) {
	response := gin.H{"error": message}
	if len(details) > 0 {
		for key, value := range details[0] {
			response[key] = value
		}
	}
	c.JSON(statusCode, response)
}
