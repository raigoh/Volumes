package utils

import (
	"log"
	"net/http"
	"text/template"
)

var templates *template.Template

func InitTemplates() error {
	var err error
	templates, err = template.ParseGlob("web/static/templates/*.html")
	return err
}

func RenderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	err := templates.ExecuteTemplate(w, tmpl, data)
	if err != nil {
		log.Printf("Error rendering template %s: %v", tmpl, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
