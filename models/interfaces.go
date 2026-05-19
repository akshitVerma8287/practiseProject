package models

import "project/dto"

type StudentProvider interface {
	GetAllStudents() ([]Student, error)
	GetStudentById(id uint) (Student, error)
	CreateStudent(student Student) (Student, error)
	UpdateStudent(id uint, req dto.UpdateStudentRequest) (Student, error)
	DeleteStudent(id uint) error
}