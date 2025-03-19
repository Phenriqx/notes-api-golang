package models

import "gorm.io/gorm"

type User struct {
    gorm.Model
	Username string `gorm:"unique"` 
	Email string `gorm:"unique"` 
	Password string `gorm:"unique"`
	Notes []Notes
}

type RegisterRequest struct {
	Username string 
	Email string
	Password string
}