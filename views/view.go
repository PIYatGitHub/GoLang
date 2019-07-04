package views

import (
	"bytes"
	"html/template"
	"io"
	"net/http"
	"path/filepath"

	"../context"
)

var (
	layoutDir   = "views/layouts/"
	templateDir = "views/"
	templateExt = ".gohtml"
)

// NewView is a function that creates and returns the new view with all the default files included
func NewView(layout string, files ...string) *View {
	addTemplatePath(files)
	addTemplateExt(files)
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

// addTemplatePath takes in a slice of strings and forms appropriate filepaths
//Eg: if you input {"home"} you will get {"views/home"}
func addTemplatePath(files []string) {
	for i, f := range files {
		files[i] = templateDir + f
	}
}

// addTemplatePath takes in a slice of strings and forms appropriate filepaths
//Eg: if you input {"views/home"} you will get {"views/home.gohtml"}
func addTemplateExt(files []string) {
	for i, f := range files {
		files[i] = f + templateExt
	}
}

func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v.Render(w, r, nil)
}

// Render is used to render the view or throw an error otherwise...
func (v *View) Render(w http.ResponseWriter, r *http.Request, data interface{}) {
	w.Header().Set("Content-type", "text/html")
	var vd Data
	switch d := data.(type) {
	case Data:
		vd = d
	default:
		vd = Data{
			Yield: data,
		}
	}
	vd.User = context.User(r.Context())
	var buff bytes.Buffer

	if err := v.Template.ExecuteTemplate(&buff, v.Layout, vd); err != nil {
		http.Error(w, "Something went wrong. Please email support@lenslocked.com if the issue persists",
			http.StatusInternalServerError)
		return
	}
	io.Copy(w, &buff)
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
