package models

import "github.com/jinzhu/gorm"

// Gallery will represent the core of our app - what the user is able to see
type Gallery struct {
	gorm.Model
	UserID uint   `gorm:"not_null;index"`
	Title  string `gorm:"not_null"`
}
