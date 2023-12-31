package main

import (
	"awesomeProject/controllers"
	"awesomeProject/database"
	"awesomeProject/middlewares"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize Database
	database.Connect("postgresql://charles:Z3SW8QMADfrxI_BdyvRmIA@matibabu-5642.8nj.cockroachlabs.cloud:26257/research?sslmode=verify-full")
	//database.Migrate()

	// Initialize Router
	router := initRouter()
	router.Run(":9000")
}

func initRouter() *gin.Engine {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	api := router.Group("/api")
	{
		api.POST("/login", controllers.GenerateToken)
		api.POST("/forgot_password", controllers.ForgotPasword)
		api.GET("/reset_password", controllers.ResetPassword)
		api.POST("/user/register", controllers.RegisterUser)
		secured := api.Group("/test").Use(middlewares.Auth())
		{
			secured.POST("/upload", controllers.Upload)
			secured.POST("/check_authentication_status", controllers.CheckAuthenticationStatus)
		}
	}
	return router
}
	