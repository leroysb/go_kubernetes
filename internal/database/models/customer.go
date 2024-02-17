package models

import "gorm.io/gorm"

type Customer struct {
	gorm.Model
	ID    int    `json:"id" gorm:"primaryKey;autoIncrement;not null"`
	Name  string `json:"name" gorm:"text;not null;default:null"`
	Email string `json:"email" gorm:"text;not null;unique"`
	Phone string `json:"phone" gorm:"text;not null;unique"`
}
