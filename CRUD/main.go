package main

import (
	"github.com/MXkodo/inventory/controllers"
	"github.com/MXkodo/inventory/initializers"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
}

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
}

func main() {
	route := gin.Default()
	route.GET("/items", controllers.GetAllItems)
	route.POST("/items", controllers.CreateItem)
	route.GET("/items/:id", controllers.GetItem)
	route.PATCH("/items/:id", controllers.UpdateItem)

	route.Run()
}