package models

import "gorm.io/gorm"

type Customer struct {
    gorm.Model
    Name  string
    Email string
    Phone string
    // Add other fields as needed
}
