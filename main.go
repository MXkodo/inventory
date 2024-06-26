package main

import (
	"net/http"

	"github.com/MXkodo/inventory/CRUD/controllers"
	"github.com/MXkodo/inventory/CRUD/initializers"
	"github.com/MXkodo/inventory/CRUD/middleware"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
	initializers.SyncItem()
	initializers.SyncAuth()
	initializers.SyncAudit()
	initializers.SyncChangeLog()
}

func main() {
	route := gin.Default()
	route.LoadHTMLGlob("public/templates/*.html")

	route.Static("/public", "./public")

	//Auth
	route.POST("/signup", controllers.Signup)
	route.POST("/login", controllers.Login)
	route.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "auth.html", gin.H{})
	})

	route.GET("/logout", controllers.Logout)

	// Главная страница
	route.GET("/", middleware.RequireAuth, func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	//CRUD
	items := route.Group("/items")
	items.Use(middleware.RequireAuth)
	{
		items.GET("", controllers.GetAllItems)
		items.POST("", controllers.CreateItem)
		items.GET("/:id", controllers.GetItem)
		items.PATCH("/:id", controllers.UpdateItem)
	}

	//Audit
	route.POST("/audit/:id", middleware.RequireAuth, controllers.InsertAudit)
	route.GET("/audit", middleware.RequireAuth, controllers.GetAllAuditItems)
	route.POST("/audit/:id/return", middleware.RequireAuth, controllers.ReturnItem)

	//ChangeLog
	route.GET("/changelog", middleware.RequireAuth, controllers.GetChangeLog)

	route.Run()
}
