package controllers

import (
	"fmt"
	"net/http"

	"../views"
)

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
	if err := r.ParseForm(); err != nil {
		panic(err)
	}
	fmt.Fprintln(w, r.PostForm["email"])
	fmt.Fprintln(w, r.PostForm["password"])
	fmt.Fprintln(w, "This is a fake passage. Pretend we created a user account...")
}
