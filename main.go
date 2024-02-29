package main

import (
	"github.com/MXkodo/inventory/CRUD/controllers"
	"github.com/MXkodo/inventory/CRUD/initializers"
	"github.com/MXkodo/inventory/CRUD/middleware"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
	initializers.SyncAuth()
}

func main() {
	route := gin.Default()

	//Auth
	route.POST("/signup", controllers.Signup)
	route.POST("/login", controllers.Login)
	route.GET("/validate", middleware.RequireAuth, controllers.Validate)

	//CRUD
	route.GET("/items", controllers.GetAllItems)
	route.POST("/items", controllers.CreateItem)
	route.GET("/items/:id", controllers.GetItem)
	route.PATCH("/items/:id", controllers.UpdateItem)

	route.Run()
}