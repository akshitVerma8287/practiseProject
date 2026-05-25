package routes

import (
	"project/controller"
	//"project/handlers"
	"project/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {

	
	admin := router.Group("/admin")
	{
		admin.POST("/create", controller.CreateAdmin)
		admin.POST("/login", controller.AdminLogin)
		
		adminProtected := admin.Group("/")
		adminProtected.Use(middleware.AuthMiddleware())
		{
			adminProtected.POST("/logout", controller.AdminLogout)
		}
	}

	studentGrp := router.Group("/api/student")
	studentGrp.Use(middleware.AuthMiddleware())

	{	

		studentGrp.GET("/getAll", controller.GetAllStudents)
		studentGrp.GET("/getById/:id", controller.GetStudentByID)
		studentGrp.POST("/create", controller.CreateStudent)
		studentGrp.PUT("/update/:id", controller.UpdateStudent)
		studentGrp.DELETE("/delete/:id", controller.DeleteStudent)

	}

	router.GET("/api/callApis", controller.FetchApis)

}