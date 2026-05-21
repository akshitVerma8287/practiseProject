package handlers

import (
	// "context"
	"log"
	"net/http"
	// "os"
	// "project/config"
	"project/models"
	"project/services"
	// "strings"


	"github.com/gin-gonic/gin"
	// "github.com/golang-jwt/jwt/v5"
)

var admins []models.Admin

func CreateAdmin(c *gin.Context) {

    var input models.Admin

    if err := c.ShouldBindJSON(&input); err != nil {

        log.Println("Invalid input:", err.Error())

        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid input",
        })

        return
    }
	
	newAdmin, err := services.CreateAdmin(input)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusCreated,newAdmin)
}

func AdminLogin(c *gin.Context) {

	var input models.Admin

	if err := c.ShouldBindJSON(&input); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input",
		})

		return
	}

	token, err := services.AdminLoginService(input)

	if err != nil {

		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func AdminLogout(c *gin.Context) {

	email := c.GetString("email")

	if email == "" {

		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})

		return
	}

	err := services.AdminLogoutService(email)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}