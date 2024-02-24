package main

import (
	"github.com/MXkodo/inventory/controllers"
	"github.com/MXkodo/inventory/models"
	"github.com/gin-gonic/gin"
)

func main() {
	route := gin.Default()
	models.ConnectDB()
	route.GET("/items", controllers.GetAllItems)
	route.POST("/items", controllers.CreateItem)

	route.Run()
}