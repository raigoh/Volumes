package errors

import (
	"literary-lions-forum/internal/utils"
	"net/http"
)

func ErrorHandler(w http.ResponseWriter, r *http.Request) {
	utils.RenderErrorTemplate(w, nil)
}
