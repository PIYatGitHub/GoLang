package views

import "html/template"

// NewView is a function that creates and returns the new view with all the default files included
func NewView(files ...string) *View {
	files = append(files, "views/layouts/footer.gohtml")
	t, err := template.ParseFiles(files...)

	if err != nil {
		panic(err)
	}

	return &View{
		Template: t,
	}

}

// View is a struct to create the template - that is all it does...
type View struct {
	Template *template.Template
}
