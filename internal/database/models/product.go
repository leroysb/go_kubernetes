package models

import "gorm.io/gorm"

type Product struct {
    gorm.Model
    Name  string
    Price float64
    Stock int
    // Add other fields as needed
}
