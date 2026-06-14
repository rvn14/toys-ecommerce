package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/rvn14/toys-ecommerce/backend/internal/response"
)

func ok(c *gin.Context, message string, data interface{}) {
	response.OK(c, message, data)
}

func created(c *gin.Context, message string, data interface{}) {
	response.Created(c, message, data)
}

func badRequest(c *gin.Context, message string, errors interface{}) {
	response.BadRequest(c, message, errors)
}

func validationError(c *gin.Context, err error) {
	response.BadRequest(c, "Validation failed", response.ValidationErrors(err))
}

func unauthorized(c *gin.Context, message string) {
	response.Unauthorized(c, message)
}

func forbidden(c *gin.Context, message string) {
	response.Forbidden(c, message)
}

func conflict(c *gin.Context, message string) {
	response.Conflict(c, message)
}

func notFound(c *gin.Context, message string) {
	response.NotFound(c, message)
}

func internalError(c *gin.Context, message string) {
	response.InternalServerError(c, message)
}
