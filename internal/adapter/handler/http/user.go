package httpserver

import (
	"1337bo4rd/internal/core/domain"
	"1337bo4rd/internal/core/port"
	"html/template"
	"net/http"
	"strings"
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
	user := getSession(r)

	if r.Method == http.MethodGet {
		data := struct {
			User *domain.User
		}{User: user}
		h.tmpl.ExecuteTemplate(w, "profile.html", data)
		return
	}

	if r.Method == http.MethodPost && r.URL.Path == "/profile/update-name" {
		newName := r.FormValue("name")
		if strings.TrimSpace(newName) == "" {
			renderError(w, h.tmpl, http.StatusBadRequest, "Name cannot be empty or whitespace")
			return
		}
		// update user and session
		user.Name = newName
		err := h.svc.UpdateUser(user)
		if err != nil {
			renderError(w, h.tmpl, statusCode, "Failed to update user")
			return
		}

		// Set updated session
		setSession(w, r, user)

		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}
}
