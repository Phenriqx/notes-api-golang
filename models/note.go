package models

import (
	"gorm.io/gorm"
)

type Notes struct {
	gorm.Model // Gorm.Model automatically generates the ID, CreatedAt, UpdatedAt and DeletedAt fields
	Title      string `gorm:"not null"`
	Content    string `gorm:"not null"`
	UserID     uint
}