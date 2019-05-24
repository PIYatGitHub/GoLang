package views

import "html/template"

// NewView is a function that creates and returns the new view with all the default files included
func NewView(layout string, files ...string) *View {
	files = append(files,
		"views/layouts/navbar.gohtml",
		"views/layouts/bootstrap.gohtml",
		"views/layouts/footer.gohtml")
	t, err := template.ParseFiles(files...)

	if err != nil {
		panic(err)
	}

	return &View{
		Template: t,
		Layout:   layout,
	}

}

// View is a struct to create the template - that is all it does...
type View struct {
	Template *template.Template
	Layout   string
}
