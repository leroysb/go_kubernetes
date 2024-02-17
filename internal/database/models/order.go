package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	Customer   Customer `gorm:"foreignKey:CustomerID"`
	CustomerID uint
	Product    Product `gorm:"foreignKey:ProductID"`
	ProductID  uint
	Quantity   int     `json:"quantity" gorm:"integer;not null;default:null"`
	Amount     float64 `json:"amount" gorm:"float;not null;default:null"`
	Time       string  `json:"time" gorm:"text;not null;default:null"`
	Status     string  `json:"status" gorm:"text;not null;default:null"`
}
