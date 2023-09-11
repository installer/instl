package templates

import (
	"embed"
	"html/template"
	"io"
)

//go:embed *.gohtml lib/*.gohtml
var tmplFS embed.FS

type Template struct {
	templates *template.Template
}

func New() *Template {
	templates := template.Must(template.New("").ParseFS(tmplFS, "*.gohtml", "lib/*.gohtml"))
	// Remove file extension from template names

	return &Template{
		templates: templates,
	}
}

func (t *Template) Load() error {
	tmpl := template.Must(t.templates.Clone())
	tmpl = template.Must(tmpl.ParseFS(tmplFS, "*.gohtml"))
	t.templates = tmpl
	return nil
}

func (t *Template) Render(w io.Writer, name string, data interface{}, s ...string) error {
	tmpl := template.Must(t.templates.Clone())
	tmpl = template.Must(tmpl.ParseFS(tmplFS, name))
	return tmpl.ExecuteTemplate(w, name, data)
}
