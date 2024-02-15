package models

import "gorm.io/gorm"

type Order struct {
    gorm.Model
    CustomerID uint
    ProductID  uint
    Quantity   int
    Amount     float64
    Time       string
    // Add other fields as needed
}
