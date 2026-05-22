package repositories

import (
	"context"
	"errors"
	"project/config"
	"project/dbops"
	"project/models"
)

func AdminAlreadyExists(email string) (bool, error) {

	query := `
	SELECT id
	FROM admins
	WHERE email = $1
	`

	rows, err := dbops.PostgresRepo.Fetch(
		query,
		email,
	)

	if err != nil {
		return false, err
	}

	defer rows.Close()

	if rows.Next() {
		return true, nil
	}

	return false, nil
}

func CreateAdmin(admin models.Admin) (models.Admin, error) {

	query := `
	INSERT INTO admins(email,password)
	VALUES($1,$2)
	RETURNING id
	`

	rows, err := dbops.PostgresRepo.Insert(
		query,
		admin.Email,
		admin.Password,
	)

	if err != nil {
		return models.Admin{}, err
	}

	defer rows.Close()

	if rows.Next() {

		err := rows.Scan(
			&admin.ID,
		)

		if err != nil {
			return models.Admin{}, err
		}
	}

	return admin, nil
}

func FindAdminByEmail(email string) (models.Admin, error) {

	query := `
	SELECT id, email, password
	FROM admins
	WHERE email = $1
	`

	rows, err := dbops.PostgresRepo.Fetch(
		query,
		email,
	)

	if err != nil {
		return models.Admin{}, err
	}

	defer rows.Close()

	var admin models.Admin

	if rows.Next() {

		err := rows.Scan(
			&admin.ID,
			&admin.Email,
			&admin.Password,
		)

		if err != nil {
			return models.Admin{}, err
		}

		return admin, nil
	}

	return models.Admin{}, errors.New("admin not found")
}

func DeleteTokenFromRedis(email string) error {

	err := config.RedisClient.Del(
		context.Background(),
		email,
	).Err()

	return err
}