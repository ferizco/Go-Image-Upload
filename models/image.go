package models

import (
	"gorm.io/gorm"
)

// Image model
type Image struct {
	gorm.Model
	ID       int
	Username string
	Subject  string
	ImageURL string
}
