package helpers

import (

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Status    bool        `json:"status"`
	Message   string      `json:"message"`
	ErrorCode string      `json:"error_code,omitempty"`
	Data      interface{} `json:"data,omitempty"`
}

func Success(c *gin.Context, statusCode int, message string, data interface{}) {

	response := APIResponse{
		Status:  true,
		Message: message,
		Data:    data,
	}

	c.JSON(statusCode, response)
}

func Error(c *gin.Context, statusCode int, message string, errorCode string) {

	response := APIResponse{
		Status:    false,
		Message:   message,
		ErrorCode: errorCode,
	}

	c.JSON(statusCode, response)
}