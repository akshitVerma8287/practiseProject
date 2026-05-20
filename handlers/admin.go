package handlers

import (
	"context"
	"log"
	"net/http"
	"os"
	"project/config"
	"project/models"
	"project/utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var admins []models.Admin

func CreateAdmin(c *gin.Context) {

    var input models.Admin
	var existingAdmin models.Admin

    if err := c.ShouldBindJSON(&input); err != nil {

        log.Println("Invalid input:", err.Error())

        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid input",
        })

        return
    }

	err := config.DB.Where(
		"email = ?",
		input.Email,
	).First(&existingAdmin).Error

	if err == nil {

		log.Println("Admin with this email already exists")
		c.JSON(http.StatusConflict, gin.H{
			"error": "Admin with this email already exists",
		})
		return 
	}
	

    hashedPassword, err := bcrypt.GenerateFromPassword(
        []byte(input.Password),
        bcrypt.DefaultCost,
    )

    if err != nil {

        log.Println("Password hashing failed:", err.Error())

        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Password hashing failed",
        })

        return
    }

    input.Password = string(hashedPassword)

    result := config.DB.Create(&input)

    if result.Error != nil {

        log.Println("Failed to create admin:", result.Error)

        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to create admin",
        })

        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Admin created successfully",
    })
}

func AdminLogin(c *gin.Context) {

	var input models.Admin
	var admin models.Admin

	// Read request body
	if err := c.ShouldBindJSON(&input); err != nil {

		log.Println("Invalid input:", err.Error())

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input",
		})

		return
	}

	// FETCH ADMIN FROM DATABASE
	err := config.DB.Where(
		"email = ?",
		input.Email,
	).First(&admin).Error

	if err != nil {

		log.Println("Admin not found:", err.Error())

		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid email or password",
		})

		return
	}

	// COMPARE HASHED PASSWORD WITH INPUT PASSWORD
	err = bcrypt.CompareHashAndPassword(
		[]byte(admin.Password), // hash from DB
		[]byte(input.Password), // plain password from request
	)

	if err != nil {

		log.Println("Password mismatch")

		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid email or password",
		})

		return
	}

	// GENERATE JWT
	token, err := utils.GenerateJWT(admin.Email)

	if err != nil {

		log.Println("Token generation failed:", err.Error())

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Token generation failed",
		})

		return
	}

	err = config.RedisClient.Set(
		context.Background(),
		input.Email,
		token,
		time.Hour*24,
	).Err()

	if err != nil {
		log.Println("Failed to store token in Redis:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to store token in Redis",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func AdminLogout(c *gin.Context) {

	// GET AUTH HEADER
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {

		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Authorization header missing",
		})

		return
	}

	// SPLIT "Bearer TOKEN"
	splitToken := strings.Split(authHeader, " ")

	if len(splitToken) != 2 {

		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid token format",
		})

		return
	}

	tokenString := splitToken[1]

	// PARSE TOKEN
	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (interface{}, error) {

			return []byte(
				os.Getenv("JWT_SECRET_KEY"),
			), nil
		},
	)

	if err != nil {
		log.Println("Token parsing failed:", err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid token",
		})

		return
	}

	// EXTRACT CLAIMS
	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		log.Println("Failed to parse claims")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid token claims",
		})

		return
	}

	// GET EMAIL
	email, ok := claims["email"].(string)

	if !ok {
		log.Println("Email not found in token claims")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Email not found in token",
		})

		return
	}

	// DELETE TOKEN FROM REDIS
	err = config.RedisClient.Del(
		context.Background(),
		email,
	).Err()

	if err != nil {

		log.Println("Failed to delete token:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Logout failed",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}