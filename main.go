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
	initializers.SyncAuth()
}

func main() {
	route := gin.Default()
	//CORS
	   route.LoadHTMLGlob("public/templates/*.html") // Загрузка шаблонов HTML

	// Public group with authentication
	route.Static("/public", "./public") // Обслуживание статических файлов из папки public

	//Auth
	route.POST("/signup", controllers.Signup)
	route.POST("/login", controllers.Login)
	route.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "auth.html", gin.H{}) // Отображение страницы авторизации
	})

	route.GET("/logout", controllers.Logout)

	// Главная страница
	route.GET("/", middleware.RequireAuth, func(c *gin.Context) {
    c.HTML(http.StatusOK, "index.html", gin.H{}) // Отображение главной страницы
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
	
	route.Run()
}
