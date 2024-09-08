package utils

import (
	"log"
	"net/http"
	"text/template"
)

// templates is a package-level variable that holds all parsed templates
var templates *template.Template

// InitTemplates initializes the templates by parsing all HTML files in the templates directory
func InitTemplates() error {
	var err error
	// ParseGlob parses all files matching the pattern and stores them in templates
	templates, err = template.ParseGlob("web/static/templates/*.html")
	return err
}

// RenderTemplate executes the specified template and writes the output to the http.ResponseWriter
func RenderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	// ExecuteTemplate applies the named template to the specified data object
	err := templates.ExecuteTemplate(w, tmpl, data)
	if err != nil {
		// Log the error for debugging purposes
		log.Printf("Error rendering template %s: %v", tmpl, err)
		// Send a generic error message to the client
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		// Note: The following line is commented out, presumably to avoid potential infinite recursion
		//RenderErrorTemplate(w, err)
		RenderErrorTemplate(w, err, http.StatusInternalServerError, "Server error, please be patient, we are doing your best :(")
	}
}

// RenderErrorTemplate renders a specific error template with the provided error
// Function handels writing the header and if some speicfic text is needed to write in the page, example "please contact support on number +358 040 1231234"
func RenderErrorTemplate(w http.ResponseWriter, err error, status int, specificText string) {
	// Create a struct to hold error data for the template
	errData := struct {
		Error        error
		SpecificText string
	}{
		Error:        err,
		SpecificText: specificText,
	}

	// If no error is provided, use a nil error
	if err == nil {
		errData.Error = err
	}

	// Write specific status for error page, 404, 400, 500, etc.
	w.WriteHeader(status)

	if status == http.StatusBadRequest {

		err := templates.ExecuteTemplate(w, "error-page-badrquest.html", errData)
		if err != nil {
			// Log the error if the error template itself fails to render
			log.Printf("Error rendering template %s: %v", "error-page.html", err)
			// Send a generic error message to the client
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			// Attempt to render the error template again with the new error
			// Note: This could potentially cause infinite recursion if not handled carefully
			RenderErrorTemplate(w, err, 400, "")
		}
	} else {
		// Execute the error template with the error data, if error is not Bad Request
		err := templates.ExecuteTemplate(w, "error-page.html", errData)
		if err != nil {
			// Log the error if the error template itself fails to render
			log.Printf("Error rendering template %s: %v", "error-page.html", err)
			// Send a generic error message to the client
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			// Attempt to render the error template again with the new error
			// Note: This could potentially cause infinite recursion if not handled carefully
			RenderErrorTemplate(w, err, 400, "")
		}
	}
}
