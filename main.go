package main

import (
	"github.com/MXkodo/inventory/controllers"
	"github.com/gin-gonic/gin"
)

func main() {
	route := gin.Default()
	route.GET("/items", controllers.GetAllItems)
	route.POST("/items", controllers.CreateItem)
	route.GET("/items/:id", controllers.GetItem)
	route.PATCH("/items/:id", controllers.UpdateItem)

	route.Run()
}