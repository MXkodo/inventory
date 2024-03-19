package controllers

import (
	"fmt"
	"net/http"
	"time"

	initial "github.com/MXkodo/inventory/CRUD/initializers"
	"github.com/MXkodo/inventory/CRUD/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreateItemInput struct {
	InvNumber string `json:"invnumber"  binding:"required"`
	Name      string `json:"name"  binding:"required"`
	Storage   string `json:"storage"  binding:"required"`
	Date      string `json:"date"`
	Budget    bool   `json:"budget"`
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
type ChangeLog struct {
    InvNumber string `json:"invnumber"`
    Name      string `json:"name"`
    Date      string `json:"date"`
    Desc      string `json:"desc"`
    Change   string `json:"change"`
    Action    string `json:"action"`
    UserName string `json:"username"`
}
//-------------------------------------Items-----------------------------------------//
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
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userName, err := ExtractUserFIO(context)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	item := models.Item{
		InvNumber: input.InvNumber,
		Name:      input.Name,
		Storage:   input.Storage,
		Date:      time.Now().Format("02.01.2006"),
		Budget:    input.Budget,
		Desc:      input.Desc}

	initial.DB.Create(&item)

	context.JSON(http.StatusOK, gin.H{"items": item})
	LogChange(initial.DB, item, "добавлено", map[string]interface{}{"Desc": item.Desc}, userName)

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

	userName, err := ExtractUserFIO(context)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
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

    changes := make(map[string]interface{})
    for key, value := range updateFields {
        changes[key] = value
    }

   LogChange(initial.DB, item, "обновлено", map[string]interface{}{"Desc": item.Desc}, userName)

   context.JSON(http.StatusOK, gin.H{"message": "Данные успешно обновлены"})

}

//----------------------------------------Audit-----------------------------------------------//

// GET /audit
// Получение всех записей из таблицы audit
func GetAllAuditItems(context *gin.Context) {
	var auditItems []models.AuditItem
	initial.DB.Find(&auditItems)
	context.JSON(http.StatusOK, gin.H{"auditItems": auditItems})
}

// POST /audit/:id
// Вставка записи в таблицу audit и удаление из таблицы item
func InsertAudit(context *gin.Context) {
	// Извлекаем ID из параметров URL
	id := context.Param("id")

	// Проверяем, существует ли элемент в таблице item
	var item models.Item
	if err := initial.DB.Where("id = ?", id).First(&item).Error; err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Элемент не найден"})
		return
	}

	userName, err := ExtractUserFIO(context)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Создаем запись для вставки в audit
	auditItem := models.AuditItem{
		InvNumber: item.InvNumber,
		Name:      item.Name,
		Storage:   item.Storage,
		Date:      time.Now().Format("02.01.2006"),
		Budget:    item.Budget,
		Desc:      item.Desc,
		DeletedAt: time.Now(),
	}

	// Сохраняем запись в таблице audit
	result := initial.DB.Create(&auditItem)
	if result.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при сохранении в audit"})
		return
	}

	// Проверяем, была ли запись успешно сохранена
	if result.RowsAffected == 0 {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Запись не была сохранена в audit"})
		return
	}

	// Удаляем запись из таблицы item
	initial.DB.Delete(&item)
	LogChange(initial.DB, item, "отправлено в архив", map[string]interface{}{"Desc": item.Desc}, userName)
	context.JSON(http.StatusOK, gin.H{})
}
// POST /audit/:id/return
// Возврат записи из таблицы audit в таблицу item
func ReturnItem(context *gin.Context) {
    // Извлекаем ID из параметров URL
    id := context.Param("id")

    // Проверяем, существует ли элемент в таблице audit
    var auditItem models.AuditItem
    if err := initial.DB.Where("id = ?", id).First(&auditItem).Error; err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"error": "Элемент не найден в audit"})
        return
    }
	userName, err := ExtractUserFIO(context)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
    // Создаем запись для возврата в таблицу item
    item := models.Item{
        InvNumber: auditItem.InvNumber,
        Name:      auditItem.Name,
        Storage:   auditItem.Storage,
        Date:      auditItem.Date,
        Budget:    auditItem.Budget,
        Desc:      auditItem.Desc,
    }

    // Сохраняем запись в таблице item
    result := initial.DB.Create(&item)
    if result.Error != nil {
        context.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при сохранении в item"})
        return
    }

    // Проверяем, была ли запись успешно сохранена
    if result.RowsAffected == 0 {
        context.JSON(http.StatusInternalServerError, gin.H{"error": "Запись не была сохранена в item"})
        return
    }

    // Удаляем запись из таблицы audit
    initial.DB.Delete(&auditItem)
	LogChange(initial.DB, item, "возвращено из архива", map[string]interface{}{"Desc": item.Desc}, userName)
    context.JSON(http.StatusOK, gin.H{"message": "Запись успешно возвращена в item и удалена из audit"})
}

//------------------------------------------ChangeLog--------------------------------------------------//

func GetChangeLog(context *gin.Context) {
    var changeLogs []models.ChangeLog
    result := initial.DB.Find(&changeLogs)
    if result.Error != nil {
        // Логирование ошибки, если она есть
        context.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return
    }
    // Логирование количества найденных записей
    fmt.Printf("Found %d change logs\n", len(changeLogs))
    context.JSON(http.StatusOK, gin.H{"changeLogs": changeLogs})
}

func LogChange(db *gorm.DB, item models.Item, action string, changes map[string]interface{}, userName string) {
    changeLog := ChangeLog{
        InvNumber: item.InvNumber,
        Name:      item.Name,
        Date:      item.Date,
        Desc:      item.Desc,
        Change:   changes["Desc"].(string),
        Action:    action,
        UserName: userName, 
    }
    if err := db.Create(&changeLog).Error; err != nil {
    // Log the error or handle it as appropriate for your application
    return
}
}

