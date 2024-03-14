package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FIO string
	UserName string `gorm:"unique"`
	Password string
}
