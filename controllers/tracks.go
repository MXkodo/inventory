package controllers

import (
	"net/http"

	"github.com/MXkodo/inventory/models"
	"github.com/gin-gonic/gin"
)

type CreateItemInput struct {
  	InvNumber string `json:"invnumber"  binding:"required"`
 	Name      string `json:"name"  binding:"required"`
 	Storage   string `json:"storage"  binding:"required"`
	Date      string `json:"date"`
 	Budget    bool   `json:"budget"  binding:"required"`
 	Desc      string `json:"desc"`
}

// GET /items
// Получаем список всех позиций
func GetAllItems(context *gin.Context) {
   db := models.ConnectDB()
   var items []models.Item
   db.Find(&items)
   context.JSON(http.StatusOK, gin.H{"items": items})
}

// POST /items
// Создание позиции
func CreateItem(context *gin.Context) {
   db := models.ConnectDB()
   var input CreateItemInput
   if err := context.ShouldBindJSON(&input); err != nil {
      context.JSON(http.StatusBadRequest, gin.H{"error":      err.Error()})
      return
   }

   item := models.Item{
	InvNumber: input.InvNumber, 
	Name: input.Name, 
	Storage: input.Storage, 
	Date: input.Date, 
	Budget: input.Budget, 
	Desc: input.Desc,}

   db.Create(&item)

   context.JSON(http.StatusOK, gin.H{"items": item})
}