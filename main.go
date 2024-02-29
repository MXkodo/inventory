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

	//CRUD
	items := route.Group("/items")
	items.Use(middleware.RequireAuth)
	{
		items.GET("", controllers.GetAllItems)
		items.POST("", controllers.CreateItem)
		items.GET("/:id", controllers.GetItem)
		items.PATCH("/:id", controllers.UpdateItem)
	}
	route.Run()
}
