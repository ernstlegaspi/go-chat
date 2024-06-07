package handlers

import (
	"gochat/internal/utils"
	"net/http"
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
	var name string = r.FormValue("name")
	var password string = r.FormValue("password")

	if email == "" || name == "" || password == "" {
		utils.AlertError(w, "All Fields are required.")
		return
	}
}
