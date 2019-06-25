package controllers

import (
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

//Galleries is the gallery struct!!!
type Galleries struct {
	New *views.View
	gs  models.GalleryService
}
