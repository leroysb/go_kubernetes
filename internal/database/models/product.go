package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	ID    int     `json:"id" gorm:"primaryKey;autoIncrement;not null"`
	Name  string  `json:"name" gorm:"text;not null;default:null"`
	Price float64 `json:"price" gorm:"float;not null;default:null"`
	Stock int     `json:"stock" gorm:"integer;not null;default:null"`
}
