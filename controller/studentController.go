package controller

import (
	"net/http"
	"strconv"

	helper "project/helper"
	"project/dto"
	"project/models"

	"github.com/gin-gonic/gin"
)

var studentProvider models.StudentProvider

func InitStudentProvider(provider models.StudentProvider) {
	studentProvider = provider
} 

// GetAllStudents godoc
// @Summary Get all students
// @Description Fetch all students 	
// @Tags Students
// @Accept json
// @Produce json
// @Success 200 {array} models.Student
// @Router /api/student/getAll [get]
func GetAllStudents(c *gin.Context) {

	students, statusCode := studentProvider.GetAllStudents()

	if statusCode != http.StatusOK {
		c.JSON(statusCode, gin.H{
			"error": "failed to fetch students",
		})
		return
	}

	c.JSON(statusCode, students)
}

// GetStudentByID godoc
// @Summary Get student by ID
// @Description Fetch student using ID
// @Tags Students
// @Accept json
// @Produce json
// @Param id path int true "Student ID"
// @Success 200 {object} models.Student
// @Router /api/student/getById/{id} [get]
func GetStudentByID(c *gin.Context) {

	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid student id",
		})
		return
	}

	student, statusCode := studentProvider.GetStudentById(uint(id))

	if statusCode != http.StatusOK {
		if statusCode == http.StatusNotFound {
			c.JSON(statusCode, gin.H{
				"error": "student not found",
			})
			return
		}
		c.JSON(statusCode, gin.H{
			"error": "failed to fetch student",
		})
		return
	}


	c.JSON(statusCode, student)
}

// CreateStudent godoc
// @Summary Create student
// @Description Create a new student
// @Tags Students
// @Accept json
// @Produce json
// @Param student body models.Student true "Student Data"
// @Success 201 {object} models.Student
// @Router /api/student/create [post]
func CreateStudent(c *gin.Context) {

	var student models.Student

	err := c.ShouldBindJSON(&student)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	newStudent, statusCode := studentProvider.CreateStudent(student)

	if statusCode != http.StatusCreated {
		if statusCode == http.StatusBadRequest {
			helper.Error(c, statusCode, "age must be greater than 16", strconv.Itoa(statusCode))
			return
		}
		helper.Error(c, statusCode, "failed to create student", strconv.Itoa(statusCode))
		return
	}

	c.JSON(statusCode, newStudent)
}

// UpdateStudent godoc
// @Summary Update student
// @Description Update student by ID
// @Tags Students
// @Accept json
// @Produce json
// @Param id path int true "Student ID"
// @Param student body dto.UpdateStudentRequest true "Updated Student"
// @Success 200 {object} models.Student
// @Router /api/student/update/{id} [put]
func UpdateStudent(c *gin.Context) {

	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid id",
		})
		return
	}

	var req dto.UpdateStudentRequest

	err = c.ShouldBindJSON(&req)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	updatedStudent, statusCode := studentProvider.UpdateStudent(uint(id), req)

	if statusCode != http.StatusOK {
		c.JSON(statusCode, gin.H{
			"error": "failed to update student",
		})
		return
	}

	c.JSON(statusCode, updatedStudent)
}

// DeleteStudent godoc
// @Summary Delete student
// @Description Delete student by ID
// @Tags Students
// @Accept json
// @Produce json
// @Param id path int true "Student ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/student/delete/{id} [delete]
func DeleteStudent(c *gin.Context) {

	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid student id",
		})
		return
	}

	statusCode := studentProvider.DeleteStudent(uint(id))

	if statusCode != http.StatusOK {
		c.JSON(statusCode, gin.H{
			"error": "failed to delete student",
		})
		return
	}

	c.JSON(statusCode, gin.H{
		"message": "student deleted successfully",
	})
}