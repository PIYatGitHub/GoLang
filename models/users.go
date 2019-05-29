package models

import "github.com/jinzhu/gorm"

// User will serve to save our users with the appropriate fields...
type User struct {
	gorm.Model
	Name  string
	Email string `gorm:"not null; unique_index"`
}
