package controllers

import (
	"fmt"
	"net/http"

	"../models"
	"../views"
)

// NewUser creates a new user view - capt. obvious strikes again!!!
// This function shall panic if there is some err.
func NewUser(us *models.UserService) *Users {
	return &Users{
		NewView:   views.NewView("bootstrap", "users/new"),
		LoginView: views.NewView("bootstrap", "users/login"),
		us:        us,
	}
}

// New  --> Use to render the form to create a new user!
// GET /signup
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	u.NewView.Render(w, nil)
}

// Login is called whenever you want to log the user in
// POST /login
func (u *Users) Login(w http.ResponseWriter, r *http.Request) {
	form := LoginForm{}
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}
	user, err := u.us.Authenticate(form.Email, form.Password)
	switch err {
	case models.ErrNotFound:
		fmt.Fprintln(w, "Invalid email address...")
	case models.ErrInvalidPass:
		fmt.Fprintln(w, "Invalid passowrd...")
	case nil:
		fmt.Fprintln(w, user)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Create is called whenever you submit the form ... se we create
// a new user account here...
// POST /signup
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	var form SignupForm
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}
	user := models.User{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	}
	if err := u.us.Create(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, user)
}

//Users is a users struct!!!
type Users struct {
	NewView   *views.View
	LoginView *views.View
	us        *models.UserService
}

// LoginForm is a struct to hold our login data, e.g. email and password
type LoginForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

// SignupForm is a struct to hold our signup data, e.g. name, email and password
type SignupForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}
