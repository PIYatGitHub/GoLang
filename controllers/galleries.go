package controllers

import (
	"fmt"
	"net/http"

	"../context"
	"../models"
	"../views"
)

// NewGallery creates a new gallery view - capt. obvious strikes again!!!
// This function shall panic if there is some err.
func NewGallery(gs models.GalleryService) *Galleries {
	return &Galleries{
		New: views.NewView("bootstrap", "galleries/new"),
		gs:  gs,
	}
}

// Create is called whenever you submit the form ... se we create
// a new user gallery here...
// POST /galleries
func (g *Galleries) Create(w http.ResponseWriter, r *http.Request) {
	var form GalleryForm
	var vd views.Data
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		g.New.Render(w, vd)
		return
	}
	user := context.User(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	fmt.Println("Create got the user correctly off the context: ", user)
	gallery := models.Gallery{
		Title:  form.Title,
		UserID: user.ID,
	}
	if err := g.gs.Create(&gallery); err != nil {
		vd.SetAlert(err)
		g.New.Render(w, vd)
		return
	}
	fmt.Fprintln(w, gallery)
}

//Galleries is the gallery struct!!!
type Galleries struct {
	New *views.View
	gs  models.GalleryService
}

// GalleryForm is a struct to hold our gallery data, e.g. the title
type GalleryForm struct {
	Title string `schema:"title"`
}
