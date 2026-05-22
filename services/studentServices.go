package services

import (
	"log"
	"net/http"
	"project/dto"
	"project/models"
	"project/repositories"
)

type StudentService struct {
}

func InitStudentService() StudentService {
	return StudentService{}
}

func (s StudentService) GetAllStudents() ([]models.Student, int) {
	students, err := repositories.GetAllStudents()
	if err != nil {
		return nil, http.StatusInternalServerError
	}
	return students, http.StatusOK
}

func (s StudentService) GetStudentById(id uint) (models.Student, int) {
	student, err := repositories.GetStudentById(id)

	
	if err != nil{
		if err.Error() == "student not found" {
			return models.Student{}, http.StatusNotFound	
		}
		return models.Student{}, http.StatusInternalServerError
	} 
	return student, http.StatusOK
}

func (s StudentService) CreateStudent(student models.Student) (models.Student, int) {

	if student.Age <= 16 {
		return models.Student{}, http.StatusBadRequest
	}
	student, err := repositories.AddStudent(student)

	if err != nil {
		return models.Student{}, http.StatusInternalServerError
	}
	
	return student, http.StatusCreated
}

func (s StudentService) UpdateStudent(id uint, req dto.UpdateStudentRequest) (models.Student, int) {

	student, err := repositories.GetStudentById(id)	

	if err != nil {
		return models.Student{}, http.StatusNotFound
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

	updatedStudent, err := repositories.UpdateStudent(student)
	if err != nil {
		
		return models.Student{}, http.StatusInternalServerError
	}
	return updatedStudent, http.StatusOK
}

func (s StudentService) DeleteStudent(id uint) int {

	_, err := repositories.GetStudentById(id)

	if err != nil {
		return http.StatusNotFound
	}

	deleteStudent := repositories.DeleteStudent(id)

	if deleteStudent != nil {
		return http.StatusInternalServerError
	}

	return http.StatusOK
}

func LogStudentsAbove18() {

	students, err := repositories.GetStudentsAbove18()

	if err != nil {
		log.Println("Error fetching students:", err)
		return
	}

	for _, student := range students {

		diff := student.Age - 18

		log.Printf(
			"ID: %d | Name: %s | Age: %d | Difference: %d\n",
			student.ID,
			student.Name,
			student.Age,
			diff,
		)
	}
}