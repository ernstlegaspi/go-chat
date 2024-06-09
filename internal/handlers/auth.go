package handlers

import (
	"database/sql"
	"fmt"
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
	hashedPassword := make(chan string, 1)
	hashError := make(chan error, 1)
	startTime := time.Now()

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

	go func() {
		// hp = hashed password
		hp, hsErr := bcrypt.GenerateFromPassword([]byte(password), 10)

		if hsErr != nil {
			hashError <- hsErr
			return
		}

		hashedPassword <- string(hp)

		close(hashedPassword)
		close(hashError)
	}()

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

	select {
	case hp := <-hashedPassword:
		var id int

		var query string = `
		insert into users (createdAt, name, email, image, password, updatedAt)
		values (NOW(), $1, $2, $3, $4, NOW())
		returning id
	`

		err := tx.QueryRow(
			query,
			email,
			"test image",
			name,
			hp,
		).Scan(&id)

		if err != nil {
			fmt.Println(err)
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

		fmt.Printf("\n Time Since: %s \n", time.Since(startTime))
		utils.AlertSuccess(201, w, "Registered successfully.")
	case he := <-hashError:
		if he != nil {
			fmt.Printf("\nError: %v\n", <-hashError)
			utils.AlertError(http.StatusInternalServerError, w, "Internal Server Error")
			return
		}
	}
}
