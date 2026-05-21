package repositories

import (
	"errors"

	"project/dbops"
	"project/models"
)

func GetAllStudents() ([]models.Student, error) {

	query := `
	SELECT id, name, email, age
	FROM students
	`

	rows, err := dbops.PostgresRepo.Fetch(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var students []models.Student

	for rows.Next() {

		var student models.Student

		err := rows.Scan(
			&student.ID,
			&student.Name,
			&student.Email,
			&student.Age,
		)

		if err != nil {
			return nil, err
		}

		students = append(
			students,
			student,
		)
	}

	return students, nil
}

func GetStudentById(id uint) (models.Student, error) {

	query := `
	SELECT id, name, email, age
	FROM students
	WHERE id = $1
	`

	rows, err := dbops.PostgresRepo.Fetch(
		query,
		id,
	)

	if err != nil {
		return models.Student{}, err
	}

	defer rows.Close()

	var student models.Student

	if rows.Next() {

		err := rows.Scan(
			&student.ID,
			&student.Name,
			&student.Email,
			&student.Age,
		)

		if err != nil {
			return models.Student{}, err
		}

		return student, nil
	}

	return models.Student{},
		errors.New("student not found")
}

func AddStudent(student models.Student) (models.Student, error) {

	query := `
	INSERT INTO students(name,email,age)
	VALUES($1,$2,$3)
	RETURNING id
	`

	rows, err := dbops.PostgresRepo.Insert(
		query,
		student.Name,
		student.Email,
		student.Age,
	)

	if err != nil {
		return models.Student{}, err
	}

	defer rows.Close()

	if rows.Next() {

		err := rows.Scan(
			&student.ID,
		)

		if err != nil {
			return models.Student{}, err
		}
	}

	return student, nil
}

func UpdateStudent(student models.Student) (models.Student, error) {

	query := `
	UPDATE students
	SET name=$1,
		email=$2,
		age=$3
	WHERE id=$4
	RETURNING id
	`

	rows, err := dbops.PostgresRepo.Update(
		query,
		student.Name,
		student.Email,
		student.Age,
		student.ID,
	)

	if err != nil {
		return models.Student{}, err
	}

	defer rows.Close()

	if rows.Next() {

		var id uint

		err := rows.Scan(&id)

		if err != nil {
			return models.Student{}, err
		}
	}

	return student, nil
}

func DeleteStudent(id uint) error {

	query := `
	DELETE FROM students
	WHERE id=$1
	RETURNING id
	`

	rows, err := dbops.PostgresRepo.Delete(
		query,
		id,
	)

	if err != nil {
		return err
	}

	defer rows.Close()

	if !rows.Next() {
		return errors.New("student not found")
	}

	return nil
}