package initializers

import "github.com/MXkodo/inventory/CRUD/models"

func SyncAudit() {
	DB.AutoMigrate(&models.AuditItem{})
}
