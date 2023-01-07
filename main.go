package main

import (
	"github.com/gin-gonic/gin"
	"go-crud-example/controller"
	"go-crud-example/initializers"
	"go-crud-example/service"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	initializers.SyncDatabase()
}
func main() {
	r := gin.Default()
	r.POST("/api/v1/user/signup", controller.Signup)
	r.POST("/api/v1/user/signin", controller.Signin)

	r.POST("/api/v1/post", service.ValidateToken, controller.Create)
	r.GET("/api/v1/post", service.ValidateToken, controller.GetAll)
	r.GET("/api/v1/post/:id", service.ValidateToken, controller.FindById)
	r.PUT("/api/v1/post/:id", service.ValidateToken, controller.Update)
	r.DELETE("/api/v1/post/:id", service.ValidateToken, controller.DeleteById)
	r.Run()
}
