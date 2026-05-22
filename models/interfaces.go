package models

import "project/dto"

type StudentProvider interface {
	GetAllStudents() ([]Student, int)
	GetStudentById(id uint) (Student, int)
	CreateStudent(student Student) (Student, int)
	UpdateStudent(id uint, req dto.UpdateStudentRequest) (Student, int)
	DeleteStudent(id uint) int
}