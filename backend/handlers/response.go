package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	CodeSuccess     = 0
	CodeBadRequest  = 400
	CodeConflict    = 409
	CodeServerError = 500
)

func success(c *gin.Context, status int, message string, data any) {
	c.JSON(status, gin.H{
		"code":         CodeSuccess,
		"message":      message,
		"errorMessage": "",
		"data":         data,
	})
}

func fail(c *gin.Context, status int, code int, message string, err error) {
	errorMessage := message
	if err != nil {
		errorMessage = err.Error()
	}

	c.JSON(status, gin.H{
		"code":         code,
		"message":      message,
		"errorMessage": errorMessage,
		"data":         nil,
	})
}

func badRequest(c *gin.Context, message string, err error) {
	fail(c, http.StatusBadRequest, CodeBadRequest, message, err)
}

func serverError(c *gin.Context, message string, err error) {
	fail(c, http.StatusInternalServerError, CodeServerError, message, err)
}
