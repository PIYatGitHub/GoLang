package views

import (
	"html/template"
	"net/http"
	"path/filepath"
)

var (
	layoutDir   = "views/layouts/"
	templateExt = ".gohtml"
)

// NewView is a function that creates and returns the new view with all the default files included
func NewView(layout string, files ...string) *View {
	files = append(files, layoutFiles()...)
	t, err := template.ParseFiles(files...)

	if err != nil {
		panic(err)
	}

	return &View{
		Template: t,
		Layout:   layout,
	}

}

// Render is used to render the view or throw an error otherwise...
func (v *View) Render(w http.ResponseWriter, data interface{}) error {
	return v.Template.ExecuteTemplate(w, v.Layout, data)
}

// View is a struct to create the template - that is all it does...
type View struct {
	Template *template.Template
	Layout   string
}

// layoutFiles returns a slice of strings to represent all layout files
func layoutFiles() []string {
	files, err := filepath.Glob(layoutDir + "*" + templateExt)
	if err != nil {
		panic(err)
	}
	return files
}
