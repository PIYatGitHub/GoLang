package models

import (
	"errors"
	"strings"

	"github.com/jinzhu/gorm"

	"../hash"
	"../rand"
	// added it not to get confused as of what is needed to run this...
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/crypto/bcrypt"
)

var (
	//ErrNotFound is returned whenever you cannot find the resource at the DB
	ErrNotFound = errors.New("models: resource not found")
	// ErrInvalidID is returned if you attempt to pass in an Id <= 0
	ErrInvalidID = errors.New("models: the ID is supposed to be greater than 0")
	// ErrInvalidPass is returned if you passed in a wrong password
	ErrInvalidPass = errors.New("models: the password provided is invalid")
)

const userPwP = "wrjg82j8#$%^&#Rweg4128y8y8suTO(24#%9ghsdbu"
const hmacSecretKey = "4wjht8wywr!^Y@$Yggwj8qeyrh139hSFYHEYFehjeo235"

// User will serve to save our users with the appropriate fields...
type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"not null; unique_index"`
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
	Remember     string `gorm:"-"`
	RememberHash string `gorm:"not null; unique_index"`
}

//UserDB will be the DB layer -->
// For all methiods:
//The methods will lookup the user by the id, email or remember token provided;
// Any of these will then return either a SINGLE user or an error.
type UserDB interface {
	//methods to perform single user query
	ByID(id uint) (*User, error)
	ByEmail(email string) (*User, error)
	ByRemember(token string) (*User, error)

	//methods to alter users
	Create(user *User) error
	Update(user *User) error
	Delete(id uint) error

	//to close the DB connection
	Close() error

	// migration helpers -- very helpful for devs
	AutoMigrate() error
	DestructiveReset() error
}

// UserService is a set of methods, serving the user CRUD.
type UserService interface {
	// Authenticate will verify the email and pass.
	//On success it will return the user, otherwise you will see an ErrRecordNotFound,
	//ErrInvalidPass or another error in general.
	Authenticate(email, password string) (*User, error)
	UserDB
}

// NewUserService will take in the newly created gorm conn and will pass it onwards...
func NewUserService(connectionInfo string) (UserService, error) {
	ug, err := newUserGorm(connectionInfo)
	if err != nil {
		return nil, err
	}
	hmac := hash.NewHMAC(hmacSecretKey)
	uv := &userValidator{
		hmac:   hmac,
		UserDB: ug,
	}
	return &userService{
		UserDB: uv,
	}, nil
}

var _ UserService = &userService{}

type userService struct {
	UserDB
}

//Authenticate will lookup the provided email and pass and will return
//a user obj for logged user and err if there isnt a user
func (us *userService) Authenticate(email, password string) (*User, error) {
	foundUser, err := us.ByEmail(email)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(foundUser.PasswordHash), []byte(password+userPwP))
	if err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			return nil, ErrInvalidPass
		default:
			return nil, err
		}
	}
	return foundUser, err
}

type userValFunc func(*User) error

func runUserValFuncs(user *User, fns ...userValFunc) error {
	for _, fn := range fns {
		if err := fn(user); err != nil {
			return err
		}
	}
	return nil
}

type userValidator struct {
	UserDB
	hmac hash.HMAC
}

//ByRemember will hash the token and will call the next layer
func (uv *userValidator) ByRemember(token string) (*User, error) {
	user := User{
		Remember: token,
	}
	if err := runUserValFuncs(&user, uv.hmacRemember); err != nil {
		return nil, err
	}
	return uv.UserDB.ByRemember(user.RememberHash)
}

//ByEmail will hash the token and will call the next layer
func (uv *userValidator) ByEmail(email string) (*User, error) {
	user := User{
		Email: email,
	}
	if err := runUserValFuncs(&user, uv.normalizeEmail); err != nil {
		return nil, err
	}
	return uv.UserDB.ByEmail(user.Email)
}

//Create here is the breakout of validation from the gorm layer
//this is why you see it calling the actual method after getting its job done
func (uv *userValidator) Create(user *User) error {
	if err := runUserValFuncs(user, uv.bcryptPassword,
		uv.defalutRemember, uv.hmacRemember); err != nil {
		return err
	}
	return uv.UserDB.Create(user)
}

//Update here is the breakout of validation from the gorm layer
//this is why you see it calling the actual method after getting its job done
func (uv *userValidator) Update(user *User) error {
	if err := runUserValFuncs(user, uv.bcryptPassword, uv.hmacRemember); err != nil {
		return err
	}
	return uv.UserDB.Update(user)
}

//delete -- again will take care of validation;
func (uv *userValidator) Delete(id uint) error {
	user := User{
		Model: gorm.Model{
			ID: id,
		},
	}
	if err := runUserValFuncs(&user, uv.idGreaterThan(0)); err != nil {
		return err
	}
	return uv.UserDB.Delete(id)
}

// it will hash a password with the predefined pepper if the pword is not
//an empty string
func (uv *userValidator) bcryptPassword(user *User) error {
	if user.Password == "" {
		return nil
	}
	pwBytes := []byte(user.Password + userPwP)
	hashedBytes, err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = ""
	return nil
}

//hmacRemember will go through the hmac of the remember token
func (uv *userValidator) hmacRemember(user *User) error {
	if user.Remember == "" {
		return nil
	}
	user.RememberHash = uv.hmac.Hash(user.Remember)
	return nil
}

//defalutRemember will add a default remember token on create if none is found...
func (uv *userValidator) defalutRemember(user *User) error {
	if user.Remember != "" {
		return nil
	}
	toekn, err := rand.RememberToken()
	if err != nil {
		return err
	}
	user.Remember = toekn
	return nil
}

// idGreaterThan is a cool little closure func which cheks your ID and returns an err if it less than 0
func (uv *userValidator) idGreaterThan(n uint) userValFunc {
	return userValFunc(func(user *User) error {
		if user.ID <= n {
			return ErrInvalidID
		}
		return nil
	})
}

//normalizeEmail will turn the email to lowerr case and trim out all the extra spaces.
func (uv *userValidator) normalizeEmail(user *User) error {
	user.Email = strings.ToLower(user.Email)
	user.Email = strings.TrimSpace(user.Email)
	return nil
}

type userGorm struct {
	db *gorm.DB
}

var _ UserDB = &userGorm{}

func newUserGorm(connectionInfo string) (*userGorm, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	return &userGorm{
		db: db,
	}, nil
}

//ByID -- userGorm version will lookup the user by id;
// it will return user,nil or nil for the user and specific user (only one)
func (ug *userGorm) ByID(id uint) (*User, error) {
	var user User
	db := ug.db.Where("id = ?", id)
	err := first(db, &user)
	return &user, err
}

//ByEmail will lookup the user by his/her email address;
// it will return user,nil or nil for the user and specific user (only one)
func (ug *userGorm) ByEmail(email string) (*User, error) {
	var user User
	db := ug.db.Where("email = ?", email)
	err := first(db, &user)
	return &user, err
}

//ByRemember will lookup the user by his/her remember token;
// it will return user,nil or nil for the user and specific user (only one)
//the method will handle the hashing for us as well
func (ug *userGorm) ByRemember(rememberHash string) (*User, error) {
	var user User
	err := first(ug.db.Where("remember_hash = ?", rememberHash), &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// first is a function to get the first match from the DB.
// DO NOT FORGET to give it a pointer on the dst object, otherwise
//you may run into major pizdec!
func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}
	return err
}

//Create does take care of creating a user or returns an error if there is sth wrong...
func (ug *userGorm) Create(user *User) error {
	return ug.db.Create(user).Error
}

//Update does take care of updating a user or returns an error if there is sth wrong...
func (ug *userGorm) Update(user *User) error {
	return ug.db.Save(user).Error
}

//Delete is a dangerous function as it deletes the user by ID. Do not use it if you are not sure...
func (ug *userGorm) Delete(id uint) error {
	user := User{Model: gorm.Model{ID: id}}
	return ug.db.Delete(&user).Error
}

// Close will terminate the connection to the DB!
func (ug *userGorm) Close() error {
	return ug.db.Close()
}

//DestructiveReset deletes the users table. NEVER EVER RUN IN PROD!!!!!
func (ug *userGorm) DestructiveReset() error {
	if err := ug.db.DropTableIfExists(&User{}).Error; err != nil {
		return err
	}
	return ug.AutoMigrate()
}

//AutoMigrate is our version of the GORM function. We will use it further down the line
func (ug *userGorm) AutoMigrate() error {
	if err := ug.db.AutoMigrate(&User{}).Error; err != nil {
		return err
	}
	return nil
}
