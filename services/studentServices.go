package services

import (
	"errors"
	"project/dto"
	"project/models"
	"project/repositories"
)

type StudentService struct {
}

func InitStudentService() StudentService {
	return StudentService{}
}

func (s StudentService) GetAllStudents() ([]models.Student, error) {
	return repositories.GetAllStudents()
}

func (s StudentService) GetStudentById(id uint) (models.Student, error) {
	
	return repositories.GetStudentById(id)			// here this 's studentService is called method receiver'
}

func (s StudentService) CreateStudent(student models.Student) (models.Student, error) {

	if student.Age <= 16 {
		return models.Student{}, errors.New("age should be greater than 16")
	}

	return repositories.AddStudent(student)
}

func (s StudentService) UpdateStudent(id uint, req dto.UpdateStudentRequest) (models.Student, error) {

	student, err := repositories.GetStudentById(id)	

	if err != nil {
		return models.Student{}, errors.New("student not found")
	}

	if req.Name != nil {
		student.Name = *req.Name
	}

	if req.Age != nil {
		student.Age = *req.Age
	}

	if req.Email != nil {
		student.Email = *req.Email
	}

	return repositories.UpdateStudent(student)
}

func (s StudentService) DeleteStudent(id uint) error {

	_, err := repositories.GetStudentById(id)

	if err != nil {
		return errors.New("student not found")
	}

	return repositories.DeleteStudent(id)
}