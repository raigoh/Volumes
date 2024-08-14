package utils

import (
	"fmt"
	"net/http"
	"path/filepath"
	"text/template"
)

var templates map[string]*template.Template

func InitTemplates() error {
	templatesDir := filepath.Join("web", "static", "templates")
	templates = make(map[string]*template.Template)

	templateFiles, err := filepath.Glob(filepath.Join(templatesDir, "*.html"))
	if err != nil {
		return fmt.Errorf("error finding template files: %v", err)
	}

	for _, file := range templateFiles {
		tmpl, err := template.ParseFiles(file)
		if err != nil {
			return fmt.Errorf("error parsing template %s: %v", file, err)
		}
		templates[filepath.Base(file)] = tmpl
	}

	return nil
}

// RenderTemplate renders a template by name
func RenderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	t, ok := templates[tmpl]
	if !ok {
		http.Error(w, fmt.Sprintf("Template %s not found", tmpl), http.StatusInternalServerError)
		return
	}

	err := t.Execute(w, data)
	if err != nil {
		http.Error(w, "Unable to render template: "+err.Error(), http.StatusInternalServerError)
	}
}
