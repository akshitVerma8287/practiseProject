package routes

import (
	"project/controller"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {

	studentGrp := router.Group("/api/student")

	{	

		studentGrp.GET("/getAll", controller.GetAllStudents)
		studentGrp.GET("/getById/:id", controller.GetStudentByID)
		studentGrp.POST("/create", controller.CreateStudent)
		studentGrp.PUT("/update/:id", controller.UpdateStudent)
		studentGrp.DELETE("/delete/:id", controller.DeleteStudent)

	}
}