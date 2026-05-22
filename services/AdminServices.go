package services

import (
	"context"
	"log"
	"net/http"
	"project/config"
	"project/models"
	"project/repositories"
	"project/utils"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func CreateAdmin(input models.Admin) (models.Admin, int) {

	// CHECK IF ADMIN ALREADY EXISTS
	adminExists, err := repositories.AdminAlreadyExists(input.Email)

	if err != nil {
		return models.Admin{}, http.StatusInternalServerError
	}

	if adminExists {
		return models.Admin{}, http.StatusConflict
	}

	// HASH PASSWORD
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(input.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {

		log.Println("Password hashing failed:", err.Error())

		return models.Admin{}, http.StatusInternalServerError
	}

	// REPLACE PLAIN PASSWORD WITH HASH
	input.Password = string(hashedPassword)

	// STORE ADMIN IN DB
	createdAdmin, err := repositories.CreateAdmin(input)

	if err != nil {

		log.Println("Failed to create admin:", err.Error())

		return models.Admin{}, http.StatusInternalServerError
	}

	return createdAdmin, http.StatusCreated
}

func AdminLoginService(input models.Admin) (string, int) {

	admin, err := repositories.FindAdminByEmail(input.Email)

	if err != nil {

		log.Println("Admin not found:", err.Error())

		return "", http.StatusNotFound
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(admin.Password),
		[]byte(input.Password),
	)

	if err != nil {

		log.Println("Password mismatch")

		return "", http.StatusUnauthorized
	}

	token, err := utils.GenerateJWT(admin.Email)

	if err != nil {

		log.Println("Token generation failed:", err.Error())

		return "", http.StatusInternalServerError
	}

	err = config.RedisClient.Set(
		context.Background(),
		input.Email,
		token,
		time.Hour*24,
	).Err()

	if err != nil {

		log.Println("Failed to store token in Redis:", err.Error())

		return "", http.StatusInternalServerError
	}

	return token, http.StatusOK
}

func AdminLogoutService(email string) int {

	err := repositories.DeleteTokenFromRedis(email)

	if err != nil {

		log.Println("Failed to delete token:", err.Error())

		return http.StatusInternalServerError
	}

	return http.StatusOK
}