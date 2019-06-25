package models

import "strings"

const (
	//ErrNotFound is returned whenever you cannot find the resource at the DB
	ErrNotFound modelError = "models: resource not found"
	// ErrInvalidID is returned if you attempt to pass in an Id <= 0
	ErrInvalidID modelError = "models: the ID is supposed to be greater than 0"
	// ErrInvalidPass is returned if you passed in a wrong password
	ErrInvalidPass modelError = "models: the password provided is invalid"
	// ErrPasswordTooShort is returned if you passed in a password which is less than 8 characters long
	ErrPasswordTooShort modelError = "models: the password must be at least 8 characters long"
	//ErrPasswordRequired is thrown whenever create i s attempted w/o password
	ErrPasswordRequired modelError = "models: you must provide a password"
	//ErrEmailRequired is returned when you pass in an empty email
	ErrEmailRequired modelError = "models: your email is required"
	//ErrInvalidEmail is returned when your email fails to match the regex
	ErrInvalidEmail modelError = "models: invalid email address"
	//ErrEmailAlreadyTaken is returned whenever the email already exists in the DB
	ErrEmailAlreadyTaken modelError = "models: this email is already taken"
	//ErrRememberTooShort is returned whenever the remember token is too short...
	ErrRememberTooShort modelError = "models: the remember token is too short"
	//ErrRememberRequired is returned whenever the remember hash is not there!
	ErrRememberRequired modelError = "models: the remember hash is required"
)

type modelError string

func (e modelError) Error() string {
	return string(e)
}

func (e modelError) Public() string {
	s := strings.Replace(string(e), "models: ", "", 1)
	return strings.Title(s)
}
