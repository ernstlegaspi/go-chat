package handlers

import (
	"gochat/internal/utils"
	"net/http"
	"strings"
)

type auth struct {
	router *http.ServeMux
}

func InitAuth(router *http.ServeMux) *auth {
	return &auth{
		router: router,
	}
}

func (a *auth) InitAuthAPI() {
	var r *http.ServeMux = a.router

	r.HandleFunc("POST /register", a.register)
}

func (a *auth) register(w http.ResponseWriter, r *http.Request) {
	var email string = r.FormValue("email")
	var name string = strings.TrimSpace(r.FormValue("name"))
	var password string = r.FormValue("password")

	if email == "" || name == "" || password == "" {
		utils.AlertError(w, "All Fields are required.")
		return
	}

	if !utils.LettersAndSpaces(name) {
		utils.AlertError(w, "Name should be letters only.")
		return
	}

	if len(name) < 4 {
		utils.AlertError(w, "Name should be 4 or more characters.")
		return
	}

	if len(email) < 7 {
		utils.AlertError(w, "Email should be 8 characters or more.")
		return
	}

	if len(password) < 7 {
		utils.AlertError(w, "Password should be 8 characters or more.")
		return
	}
}
