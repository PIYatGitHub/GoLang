package controllers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"../context"
	"../models"
	"../views"
	"github.com/gorilla/mux"
)

const (
	maxMultipartMem = 30 << 20 //equals 30MB!
)

// NewGallery creates a new gallery view - capt. obvious strikes again!!!
// This function shall panic if there is some err.
func NewGallery(gs models.GalleryService, r *mux.Router) *Galleries {
	return &Galleries{
		New:       views.NewView("bootstrap", "galleries/new"),
		ShowView:  views.NewView("bootstrap", "galleries/show"),
		EditView:  views.NewView("bootstrap", "galleries/edit"),
		IndexView: views.NewView("bootstrap", "galleries/index"),
		gs:        gs,
		r:         r,
	}
}

// Index gets every single gallery the user has
// GET /galleries
func (g *Galleries) Index(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())
	galleries, err := g.gs.ByUserID(user.ID)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	var vd views.Data
	vd.Yield = galleries
	g.IndexView.Render(w, r, vd)
}

// Show is called whenever you fetch all your galleries
// GET /galleries/:id
func (g *Galleries) Show(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}
	var vd views.Data
	vd.Yield = gallery
	g.ShowView.Render(w, r, vd)
}

// Edit is called whenever you want to edit your gallery
// GET /galleries/:id/edit
func (g *Galleries) Edit(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}
	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "Gallery not found", http.StatusNotFound)
		return
	}
	var vd views.Data
	vd.Yield = gallery
	g.EditView.Render(w, r, vd)
}

// Update is called whenever you want to submit the updates to your gallery
// POST /galleries/:id/update
func (g *Galleries) Update(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}
	var form GalleryForm
	var vd views.Data
	vd.Yield = gallery

	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "Gallery not found!", http.StatusNotFound)
		return
	}

	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		g.EditView.Render(w, r, vd)
		return
	}

	gallery.Title = form.Title
	err = g.gs.Update(gallery)
	if err != nil {
		vd.SetAlert(err)
		g.EditView.Render(w, r, vd)
		return
	}
	vd.Alert = &views.Alert{
		Level:   views.AlertLvlSuccess,
		Message: "Gallery Successfully Updated!",
	}
	g.EditView.Render(w, r, vd)
}

// Create is called whenever you submit the form ... se we create
// a new user gallery here...
// POST /galleries
func (g *Galleries) Create(w http.ResponseWriter, r *http.Request) {
	var form GalleryForm
	var vd views.Data
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		g.New.Render(w, r, vd)
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
		g.New.Render(w, r, vd)
		return
	}
	url, err := g.r.Get("edit_gallery").URL("id", fmt.Sprintf("%v", gallery.ID))
	if err != nil {
		// TODO: Make this go to the index page
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	http.Redirect(w, r, url.Path, http.StatusFound)
}

// ImageUpload is called whenever you want to upload some images to your gallery
// POST /galleries/:id/images
func (g *Galleries) ImageUpload(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}
	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "Gallery not found!", http.StatusNotFound)
		return
	}
	// TODO: parse a multipart form
	var vd views.Data
	vd.Yield = gallery
	err = r.ParseMultipartForm(maxMultipartMem)
	if err != nil {
		vd.SetAlert(err)
		g.EditView.Render(w, r, vd)
		return
	}
	//create the dir to contain our images
	galleryPath := fmt.Sprintf("images/galleries/%v/", gallery.ID)
	err = os.MkdirAll(galleryPath, 0755) // hackerage...
	if err != nil {
		vd.SetAlert(err)
		g.EditView.Render(w, r, vd)
		return
	}

	files := r.MultipartForm.File["images"] //could err if the field is NOT named images
	for _, f := range files {
		file, err := f.Open()
		if err != nil {
			vd.SetAlert(err)
			g.EditView.Render(w, r, vd)
			return
		}
		defer file.Close() // file is now open and will close eventually

		dst, err := os.Create(galleryPath + f.Filename)
		if err != nil {
			vd.SetAlert(err)
			g.EditView.Render(w, r, vd)
			return
		}
		defer dst.Close()
		_, err = io.Copy(dst, file)
		if err != nil {
			vd.SetAlert(err)
			g.EditView.Render(w, r, vd)
			return
		}
	}
	fmt.Fprintln(w, "Files uploaded succesfully...")
}

// Delete is called whenever you want to delete your gallery
// POST /galleries/:id/delete
func (g *Galleries) Delete(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}
	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "Gallery not found!", http.StatusNotFound)
		return
	}
	var vd views.Data
	err = g.gs.Delete(gallery.ID)
	if err != nil {
		vd.SetAlert(err)
		vd.Yield = gallery
		g.EditView.Render(w, r, vd)
		return
	}
	http.Redirect(w, r, "/galleries", http.StatusFound)
}

// helper to get galleries by id
func (g *Galleries) galleryByID(w http.ResponseWriter, r *http.Request) (*models.Gallery, error) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid gallery ID", http.StatusNotFound)
		return nil, err
	}
	gallery, err := g.gs.ByID(uint(id))
	if err != nil {
		switch err {
		case models.ErrNotFound:
			http.Error(w, "Gallery not found", http.StatusNotFound)
		default:
			http.Error(w, "Whoops! Something went very wrong...", http.StatusInternalServerError)
		}
		return nil, err
	}
	return gallery, nil
}

//Galleries is the gallery struct!!!
type Galleries struct {
	New       *views.View
	ShowView  *views.View
	EditView  *views.View
	IndexView *views.View
	gs        models.GalleryService
	r         *mux.Router
}

// GalleryForm is a struct to hold our gallery data, e.g. the title
type GalleryForm struct {
	Title string `schema:"title"`
}
