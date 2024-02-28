package controllers

import (
	"net/http"

	initial "github.com/MXkodo/inventory/initializers"
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
type UpdateItemInput struct {
   InvNumber string `json:"invnumber"`
 	Name      string `json:"name"`
 	Storage   string `json:"storage"`
	Date      string `json:"date"`
 	Budget    bool   `json:"budget"`
 	Desc      string `json:"desc"`
}

// GET /items
// Получаем список всех позиций
func GetAllItems(context *gin.Context) {
   var items []models.Item
   initial.DB.Find(&items)
   context.JSON(http.StatusOK, gin.H{"items": items})
}

// POST /items
// Создание позиции
func CreateItem(context *gin.Context) {
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

   initial.DB.Create(&item)

   context.JSON(http.StatusOK, gin.H{"items": item})
}

// GET /items/:id
// Получение одной позиции по ID
func GetItem(context *gin.Context) {
   // Проверяем имеется ли запись
   var item models.Item
   if err := initial.DB.Where("id = ?", context.Param("id")).First(&item).Error; err != nil {
      context.JSON(http.StatusBadRequest, gin.H{"error": "Запись не существует"})
      return
   }

   context.JSON(http.StatusOK, gin.H{"tracks": item})
}

// PATCH /items/:id
// Изменения информации
func UpdateItem(context *gin.Context) {
   // Проверяем имеется ли такая запись перед тем как её менять
   var item models.Item
   if err := initial.DB.Where("id = ?", context.Param("id")).First(&item).Error; err != nil {
      context.JSON(http.StatusBadRequest, gin.H{"error": "Запись не существует"})
      return
   }

   var input UpdateItemInput
   if err := context.ShouldBindJSON(&input); err != nil {
      context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
      return
   }

	updateFields := map[string]interface{}{
		"InvNumber": input.InvNumber,
		"Name":      input.Name,
		"Storage":   input.Storage,
		"Date":      input.Date,
		"Budget":    input.Budget,
		"Desc":      input.Desc,
	}

   initial.DB.Model(&item).Updates(updateFields)

   context.JSON(http.StatusOK, gin.H{"items": item})
}