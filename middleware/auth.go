package middleware

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	"project/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {

		// GET AUTHORIZATION HEADER
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {

			log.Println("Authorization header missing")

			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header missing",
			})

			c.Abort()
			return
		}

		// SPLIT "Bearer TOKEN"
		splitToken := strings.Split(authHeader, " ")

		if len(splitToken) != 2 {

			log.Println("Invalid token format")

			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token format",
			})

			c.Abort()
			return
		}

		// ACTUAL JWT TOKEN
		tokenString := splitToken[1]

		// PARSE JWT TOKEN
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
				"error": "Token parsing failed",
			})

			c.Abort()
			return
		}

		// CHECK TOKEN VALIDITY
		if !token.Valid {

			log.Println("Invalid token")

			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
			})

			c.Abort()
			return
		}

		// EXTRACT CLAIMS
		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok {

			log.Println("Failed to parse claims")

			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token claims",
			})

			c.Abort()
			return
		}

		// EXTRACT EMAIL FROM JWT
		email, ok := claims["email"].(string)

		if !ok {

			log.Println("Email not found in token")

			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token payload",
			})

			c.Abort()
			return
		}

		c.Set("email", email)

		// FETCH TOKEN FROM REDIS
		storedToken, err := config.RedisClient.Get(
			context.Background(),
			email,
		).Result()

		if err != nil {

			log.Println("Token not found in Redis:", err.Error())

			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Session expired or logged out",
			})

			c.Abort()
			return
		}

		// COMPARE REDIS TOKEN VS CLIENT TOKEN
		if storedToken != tokenString {

			log.Println("Token mismatch")

			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid session",
			})

			c.Abort()
			return
		}

		// ALLOW REQUEST
		c.Next()
	}
}