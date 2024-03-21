package initializers

import "github.com/MXkodo/inventory/CRUD/models"

func SyncAuth() {
	DB.AutoMigrate(&models.User{})
}

