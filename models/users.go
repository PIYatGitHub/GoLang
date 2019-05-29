package models

import "github.com/jinzhu/gorm"

// UserService holds the logic?
type UserService struct {
}

//ByID will lookup the user by id;
// it will return user,nil or nil for the user and specific user (only one)
func (us *UserService) ByID(id uint) (*User, error) {

}

// User will serve to save our users with the appropriate fields...
type User struct {
	gorm.Model
	Name  string
	Email string `gorm:"not null; unique_index"`
}
