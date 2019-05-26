package controllers

import (
	"net/http"

	"../views"
)

// NewUser creates a new user - capt. obvious strikes again!!!
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
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	u.NewView.Render(w, nil)
}
