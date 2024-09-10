package errors

import (
	"html/template"
	"log"
	"net/http"
)

func RecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Log the error (optional)
				log.Printf("Recovered from panic: %v", err)

				// Set the status code for the error page
				w.WriteHeader(http.StatusInternalServerError)

				// Render your error page template with status and specific error message
				tmpl, _ := template.ParseFiles("error-page.html") // Load your error page template
				tmpl.Execute(w, map[string]interface{}{
					"Status":       500,
					"SpecificText": "Oops! Something went wrong.",
					"Error":        err,
				})
			}
		}()
		next.ServeHTTP(w, r)
	})
}
