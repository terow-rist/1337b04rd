package httpserver

import (
	"1337bo4rd/internal/core/domain"
	"1337bo4rd/internal/core/port"
	"html/template"
	"net/http"
)

type UserHandler struct {
	svc  port.UserService
	tmpl *template.Template
}

func NewUserHandler(svc port.UserService) *UserHandler {
	tmpl := template.Must(template.ParseGlob("templates/*.html"))

	return &UserHandler{
		svc:  svc,
		tmpl: tmpl,
	}
}

func (h *UserHandler) HandleProfile(w http.ResponseWriter, r *http.Request) {
	data := struct {
		User *domain.User
	}{
		User: getSession(r),
	}

	if r.Method == http.MethodGet {
		h.tmpl.ExecuteTemplate(w, "profile.html", data)
		return
	}
}
