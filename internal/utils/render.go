package utils

import (
	"literary-lions-forum/internal/models"
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

func RenderTemplate(w http.ResponseWriter, tmpl string, data models.PageData) {
	// log.Printf("Rendering template: %s with data: %+v", tmpl, data)
	if templates == nil {
		// log.Println("Templates not initialized")
		http.Error(w, "Templates not initialized", http.StatusInternalServerError)
		return
	}
	err := templates.ExecuteTemplate(w, tmpl, data)
	if err != nil {
		// log.Printf("Error executing template %s: %v", tmpl, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		// log.Printf("Template %s rendered successfully", tmpl)
	}
}
