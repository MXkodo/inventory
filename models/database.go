package models

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
    var err error
    DB, err = gorm.Open(postgres.Open(ConnectDB), &gorm.Config{})
    if err != nil {
        panic("Failed to connect to database: " + err.Error())
    }
}