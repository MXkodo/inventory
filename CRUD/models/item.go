package models

type Item struct {
	ID        uint   `json:"id" gorm:"primary_key"`
	InvNumber string `json:"invnumber" gorm:"not null"`
	Name      string `json:"name" gorm:"not null"`
	Storage   string `json:"storage" gorm:"not null"`
	Date      string `json:"date"`
	Budget    bool   `json:"budget" gorm:"not null"`
	Desc      string `json:"desc"`
}