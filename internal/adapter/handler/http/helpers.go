package httpserver

import (
	"html/template"
	"net/http"
)

func renderError(w http.ResponseWriter, tmpl *template.Template, status int, message string) {
	w.WriteHeader(status)
	tmpl.ExecuteTemplate(w, "error.html", struct {
		StatusCode int
		Message    string
	}{
		StatusCode: status,
		Message:    message,
	})
}
