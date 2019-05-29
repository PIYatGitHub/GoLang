package models

import (
	"errors"

	"github.com/jinzhu/gorm"
	// added it not to get confused as of what is needed to run this...
	_ "github.com/jinzhu/gorm/dialects/postgres"
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

//ByEmail will lookup the user by his/her email address;
// it will return user,nil or nil for the user and specific user (only one)
func (us *UserService) ByEmail(email string) (*User, error) {
	var user User
	err := us.db.Where("email = ?", email).First(&user).Error
	switch err {
	case nil:
		return &user, nil
	case gorm.ErrRecordNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

//Create does take care of creating a user or returns an error if there is sth wrong...
func (us *UserService) Create(user *User) error {
	return us.db.Create(user).Error
}

//Update does take care of updating a user or returns an error if there is sth wrong...
func (us *UserService) Update(user *User) error {
	return us.db.Save(user).Error
}

// Close will terminate the connection to the DB!
func (us *UserService) Close() error {
	return us.db.Close()
}

//DestructiveReset deletes the users table. NEVER EVER RUN IN PROD!!!!!
func (us *UserService) DestructiveReset() {
	us.db.DropTableIfExists(&User{})
	us.db.AutoMigrate(&User{})
}

// User will serve to save our users with the appropriate fields...
type User struct {
	gorm.Model
	Name  string
	Email string `gorm:"not null; unique_index"`
}
