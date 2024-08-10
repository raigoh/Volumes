package utils

import (
	"net/http"
	"path/filepath"
	"text/template"
)

var templates *template.Template

func InitTemplates() error {
	templatesDir := filepath.Join("web", "static", "templates")
	var err error
	templates, err = template.ParseGlob(filepath.Join(templatesDir, "*.html"))
	return err
}

// RenderTemplate renders a template by name
func RenderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	err := templates.ExecuteTemplate(w, tmpl, data)
	if err != nil {
		http.Error(w, "Unable to render template", http.StatusInternalServerError)
	}
}
