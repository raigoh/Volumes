package errors

import (
	"literary-lions-forum/internal/utils"
	"net/http"
)

// ErrorHandler is an HTTP handler function that renders a generic error page.
// This function can be used to handle various types of errors in a consistent manner.
// This handler is for testing porpuses
func ErrorHandler(w http.ResponseWriter, r *http.Request) {
	// Call the RenderErrorTemplate function from the utils package.
	// This function is responsible for rendering the error template.
	// The 'nil' argument suggests that no specific error details are being passed,
	// which means this might render a generic error page.
	utils.RenderErrorTemplate(w, nil, http.StatusInternalServerError, "Have you tried turning it off and on again?")
}

// ErrorHandler is an HTTP handler function that renders a generic error page.
// This function can be used to handle various types of errors in a consistent manner.
// This handler is for testing porpuses
func ErrorHandler2(w http.ResponseWriter, r *http.Request) {
	// Call the RenderErrorTemplate function from the utils package.
	// This function is responsible for rendering the error template.
	// The 'nil' argument suggests that no specific error details are being passed,
	// which means this might render a generic error page.
	utils.RenderErrorTemplate(w, nil, http.StatusBadRequest, "Have you tried turning it off and on again?")
}
