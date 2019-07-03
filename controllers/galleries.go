package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"../context"
	"../models"
	"../views"
	"github.com/gorilla/mux"
)

// NewGallery creates a new gallery view - capt. obvious strikes again!!!
// This function shall panic if there is some err.
func NewGallery(gs models.GalleryService, r *mux.Router) *Galleries {
	return &Galleries{
		New:      views.NewView("bootstrap", "galleries/new"),
		ShowView: views.NewView("bootstrap", "galleries/show"),
		gs:       gs,
		r:        r,
	}
}

// Show is called whenever you fetch all your galleries
// GET /galleries/:id
func (g *Galleries) Show(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var vd views.Data
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid gallery ID", http.StatusNotFound)
		return
	}
	gallery, err := g.gs.ByID(uint(id))
	if err != nil {
		switch err {
		case models.ErrNotFound:
			http.Error(w, "Gallery not found", http.StatusNotFound)
		default:
			http.Error(w, "Whoops! Something went very wrong...", http.StatusInternalServerError)
		}
		return
	}
	vd.Yield = gallery
	g.ShowView.Render(w, vd)
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
	url, err := g.r.Get("show_gallery").URL("id", fmt.Sprintf("%v", gallery.ID))
	if err != nil {
		// TODO: Make this go to the index page
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	http.Redirect(w, r, url.Path, http.StatusFound)
}

//Galleries is the gallery struct!!!
type Galleries struct {
	New      *views.View
	ShowView *views.View
	gs       models.GalleryService
	r        *mux.Router
}

// GalleryForm is a struct to hold our gallery data, e.g. the title
type GalleryForm struct {
	Title string `schema:"title"`
}
