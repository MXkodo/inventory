package initializers

import "github.com/MXkodo/inventory/CRUD/models"

func SyncItem() {
	DB.AutoMigrate(&models.Item{})
}
