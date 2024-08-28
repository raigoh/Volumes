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
		//RenderErrorTemplate(w, err)
	}
}

func RenderErrorTemplate(w http.ResponseWriter, err error) {
	errData := struct {
		Error        error
		SpecificText string
	}{
		Error:        err,
		SpecificText: "Have you tried turning it off and back on again ?",
	}
	if err == nil {
		errData.Error = err
	}
	err2 := templates.ExecuteTemplate(w, "error-page.html", errData)
	if err2 != nil {
		log.Printf("Error rendering template %s: %v", "error-page.html", err2)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		RenderErrorTemplate(w, err2)
	}
}
