package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"project/config"
	"project/models"
	"project/repositories"
	"project/utils"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func CreateAdmin(input models.Admin) (models.Admin, error) {

	// CHECK IF ADMIN ALREADY EXISTS
	adminExists, err := repositories.AdminAlreadyExists(input.Email)

	if err != nil {
		return models.Admin{}, err
	}

	if adminExists {
		return models.Admin{}, fmt.Errorf(
			"admin with email %s already exists",
			input.Email,
		)
	}

	// HASH PASSWORD
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(input.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {

		log.Println("Password hashing failed:", err.Error())

		return models.Admin{}, err
	}

	// REPLACE PLAIN PASSWORD WITH HASH
	input.Password = string(hashedPassword)

	// STORE ADMIN IN DB
	createdAdmin, err := repositories.CreateAdmin(input)

	if err != nil {

		log.Println("Failed to create admin:", err.Error())

		return models.Admin{}, err
	}

	return createdAdmin, nil
}

func AdminLoginService(input models.Admin) (string, error) {

	admin, err := repositories.FindAdminByEmail(input.Email)

	if err != nil {

		log.Println("Admin not found:", err.Error())

		return "", errors.New("invalid email")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(admin.Password),
		[]byte(input.Password),
	)

	if err != nil {

		log.Println("Password mismatch")

		return "", errors.New("invalid email password")
	}

	token, err := utils.GenerateJWT(admin.Email)

	if err != nil {

		log.Println("Token generation failed:", err.Error())

		return "", errors.New("token generation failed")
	}

	err = config.RedisClient.Set(
		context.Background(),
		input.Email,
		token,
		time.Hour*24,
	).Err()

	if err != nil {

		log.Println("Failed to store token in Redis:", err.Error())

		return "", errors.New("failed to store token in redis")
	}

	return token, nil
}

func AdminLogoutService(email string) error {

	err := repositories.DeleteTokenFromRedis(email)

	if err != nil {

		log.Println("Failed to delete token:", err.Error())

		return errors.New("logout failed")
	}

	return nil
}