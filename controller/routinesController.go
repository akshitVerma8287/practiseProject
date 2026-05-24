package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"project/services"
	"project/utils"
)

func FetchApis(c *gin.Context) {

	utils.ClearLogs()

	services.CallMultipleApis()

	c.JSON(
		http.StatusOK,
		gin.H{
			"message": "Completed",
			"logs":    utils.GetLogs(),
		},
	)
}