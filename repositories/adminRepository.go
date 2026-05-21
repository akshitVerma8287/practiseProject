package repositories

import (
	"context"
	"log"
	"project/config"
	"project/models"

	"gorm.io/gorm"
)

func AdminAlreadyExists(email string) (bool, error) {

	var existingAdmin models.Admin

	err := config.DB.Where(
		"email = ?",
		email,
	).First(&existingAdmin).Error

	// ADMIN FOUND
	if err == nil {

		log.Println("Admin with this email already exists")

		return true, nil
	}

	// SOME OTHER DATABASE ERROR
	if err != gorm.ErrRecordNotFound {

		log.Println("Database error:", err.Error())

		return false, err
	}

	// ADMIN DOES NOT EXIST
	return false, nil
}

func CreateAdmin(input models.Admin) (models.Admin, error) {
	result := config.DB.Create(&input)

	return input, result.Error
}

func FindAdminByEmail(email string) (models.Admin, error) {

	var admin models.Admin
	err := config.DB.Where(
		"email = ?",
		email,
	).First(&admin).Error	
	return admin, err
}

func DeleteTokenFromRedis(email string) error {

	err := config.RedisClient.Del(
		context.Background(),
		email,
	).Err()

	return err
}