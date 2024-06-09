package handlers

import (
	"database/sql"
	"gochat/internal/utils"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type auth struct {
	db     *sql.DB
	router *http.ServeMux
}

func InitAuth(db *sql.DB, router *http.ServeMux) *auth {
	return &auth{
		db:     db,
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
		utils.AlertError(http.StatusBadRequest, w, "All Fields are required.")
		return
	}

	if !utils.LettersAndSpaces(name) {
		utils.AlertError(http.StatusBadRequest, w, "Name should be letters only.")
		return
	}

	if len(name) < 4 {
		utils.AlertError(http.StatusBadRequest, w, "Name should be 4 or more characters.")
		return
	}

	if len(email) < 7 {
		utils.AlertError(http.StatusBadRequest, w, "Email should be 8 characters or more.")
		return
	}

	if len(password) < 7 {
		utils.AlertError(http.StatusBadRequest, w, "Password should be 8 characters or more.")
		return
	}

	hashedPassword, hashErr := bcrypt.GenerateFromPassword([]byte(password), 14)

	if hashErr != nil {
		utils.AlertError(http.StatusInternalServerError, w, "Internal Server Error")
		return
	}

	tx, txErr := a.db.Begin()

	if txErr != nil {
		utils.AlertError(http.StatusInternalServerError, w, "Internal Server Error")
		return
	}

	defer func() {
		if rb := recover(); rb != nil {
			tx.Rollback()
			panic(rb)
		} else if txErr != nil {
			tx.Rollback()
			utils.AlertError(http.StatusInternalServerError, w, "Internal Server Error")
			return
		} else {
			if err := tx.Commit(); err != nil {
				utils.AlertError(http.StatusInternalServerError, w, "Internal Server Error")
			}
		}
	}()

	var id int

	var query string = `
		insert into users (createdAt, name, email, password, updatedAt)
		values (NOW, $1, $2, $3, NOW())
		returning id
	`

	err := tx.QueryRow(
		query,
		name,
		email,
		hashedPassword,
	).Scan(&id)

	if err != nil {
		utils.AlertError(http.StatusInternalServerError, w, "Internal Server Error")
		return
	}

	token := utils.CreateJWT(id)

	cookie := &http.Cookie{
		Name:     "session_token",
		HttpOnly: true,
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	}

	http.SetCookie(w, cookie)
}
