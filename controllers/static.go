package controllers

import (
	"lenslocked.com/views"
)

// NewStatic instantiates the static pages as views...
func NewStatic() *Static {
	return &Static{
		Home:    views.NewView("bootstrap", "static/home"),
		Contact: views.NewView("bootstrap", "static/contact"),
	}
}

//Static is a s truct to hold all your static views - it may contain as amay as you will need...
type Static struct {
	Home    *views.View
	Contact *views.View
}
