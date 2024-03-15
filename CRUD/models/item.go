package models

import "time"

type Item struct {
	ID        uint   `json:"id" gorm:"primary_key"`
	InvNumber string `json:"invnumber" gorm:"not null"`
	Name      string `json:"name" gorm:"not null"`
	Storage   string `json:"storage" gorm:"not null"`
	Date      string `json:"date"`
	Budget    bool   `json:"budget" gorm:"not null"`
	Desc      string `json:"desc"`
}

type AuditItem struct {
	ID        uint      `gorm:"primary_key"`
	InvNumber string    `json:"invnumber"`
	Name      string    `json:"name"`
	Storage   string    `json:"storage"`
	Date      string 	`json:"date"`
	Budget    bool      `json:"budget"`
	Desc      string    `json:"desc"`
	DeletedAt time.Time `json:"deleted_at"`
}