package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"			
)

func GenerateJWT(email string) (string, error) {

	secretKey := []byte(os.Getenv("JWT_SECRET_KEY"))	// Converting to bytes bcz signing fns needs bytes

	claims := jwt.MapClaims{
		"email": email,
		"role":  "admin",
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	}														// This is the payload of token

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,			// This tells jwt to insert the header using HS256 algo
		claims,
	)

	return token.SignedString(secretKey)
}