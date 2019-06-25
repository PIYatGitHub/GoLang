package views

import (
	"html/template"
	"net/http"
	"path/filepath"
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
	if err := v.Render(w, nil); err != nil {
		panic(err)
	}
}

// Render is used to render the view or throw an error otherwise...
func (v *View) Render(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-type", "text/html")
	switch data.(type) {
	case Data: //do nothing
	default:
		data = Data{
			Yield: data,
		}

	}
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
