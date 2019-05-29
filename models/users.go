package models

import (
	"errors"

	"github.com/jinzhu/gorm"
)

var (
	//ErrNotFound is returned whenever you cannot find the resource at the DB
	ErrNotFound = errors.New("models:resource not found")
)

// NewUserService will open a singular connection to the DB!
func NewUserService(connectionInfo string) (*UserService, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	return &UserService{
		db: db,
	}, nil
}

// UserService holds the logic?
type UserService struct {
	db *gorm.DB
}

//ByID will lookup the user by id;
// it will return user,nil or nil for the user and specific user (only one)
func (us *UserService) ByID(id uint) (*User, error) {
	var user User
	err := us.db.Where("id = ?", id).First(&user).Error
	switch err {
	case nil:
		return &user, nil
	case gorm.ErrRecordNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// Close will terminate the connection to the DB!
func (us *UserService) Close() error {
	return us.db.Close()
}

// User will serve to save our users with the appropriate fields...
type User struct {
	gorm.Model
	Name  string
	Email string `gorm:"not null; unique_index"`
}
