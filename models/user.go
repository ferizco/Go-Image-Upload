package models

import (
	"gorm.io/gorm"
)

// User model
type User struct {
	gorm.Model
	Username string `gorm:"unique"`
	Password string
	Email    string
}
