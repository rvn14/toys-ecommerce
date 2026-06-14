package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

func JSON(c *gin.Context, statusCode int, success bool, message string, data interface{}, errors interface{}) {
	c.JSON(statusCode, APIResponse{
		Success: success,
		Message: message,
		Data:    data,
		Errors:  errors,
	})
}

func OK(c *gin.Context, message string, data interface{}) {
	JSON(c, http.StatusOK, true, message, data, nil)
}

func Created(c *gin.Context, message string, data interface{}) {
	JSON(c, http.StatusCreated, true, message, data, nil)
}

func BadRequest(c *gin.Context, message string, errors interface{}) {
	JSON(c, http.StatusBadRequest, false, message, nil, errors)
}

func Unauthorized(c *gin.Context, message string) {
	JSON(c, http.StatusUnauthorized, false, message, nil, nil)
}

func Forbidden(c *gin.Context, message string) {
	JSON(c, http.StatusForbidden, false, message, nil, nil)
}

func NotFound(c *gin.Context, message string) {
	JSON(c, http.StatusNotFound, false, message, nil, nil)
}

func Conflict(c *gin.Context, message string) {
	JSON(c, http.StatusConflict, false, message, nil, nil)
}

func InternalServerError(c *gin.Context, message string) {
	JSON(c, http.StatusInternalServerError, false, message, nil, nil)
}
