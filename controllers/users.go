package controllers

import (
	"fmt"
	"net/http"

	"../views"
)

// SignupForm is a struct to hold our data, e.g. email and password
type SignupForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

// NewUser creates a new user view - capt. obvious strikes again!!!
// This function shall panic if there is some err.
func NewUser() *Users {
	return &Users{
		NewView: views.NewView("bootstrap", "views/users/new.gohtml"),
	}
}

//Users is a users struct!!!
type Users struct {
	NewView *views.View
}

// New  --> Use to render the form to create a new user!
// GET /signup
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	u.NewView.Render(w, nil)
}

// Create is called whenever you submit the form ... se we create
// a new user account here...
// POST /signup
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	var form SignupForm
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}
	fmt.Fprintln(w, form)
}
