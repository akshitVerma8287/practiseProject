package controller

import (
	"log"
	"net/http"
	"project/models"
	"project/services"
	"github.com/gin-gonic/gin"
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
	
	newAdmin, statusCode := services.CreateAdmin(input)

	if statusCode != http.StatusCreated {
		c.JSON(statusCode, gin.H{
			"error": "Failed to create admin",
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

	token, statusCode := services.AdminLoginService(input)

	if statusCode != http.StatusOK {

		c.JSON(statusCode, gin.H{
			"error": "Invalid email or password",
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

	statusCode := services.AdminLogoutService(email)

	if statusCode != http.StatusOK {

		c.JSON(statusCode, gin.H{
			"error": "Failed to logout",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}