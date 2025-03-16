package models

import "gorm.io/gorm"

type User struct {
    gorm.Model
	username string `gorm:"unique"` 
	Email string `gorm:"unique"` 
	Password string `gorm:"unique"`
	Notes []Notes
}