package initializers

import "github.com/MXkodo/inventory/CRUD/models"

func SyncChangeLog() {
	DB.AutoMigrate(&models.ChangeLog{})
}
