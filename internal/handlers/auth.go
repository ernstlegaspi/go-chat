package handlers

import (
	"database/sql"
	"fmt"
	"gochat/internal/utils"
	"net/http"
	"strings"
	"time"

	"gochat/internal/models"

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

	r.HandleFunc("POST /login", a.login)
	r.HandleFunc("POST /logout", a.logout)
	r.HandleFunc("POST /register", a.register)
}

func (a *auth) login(w http.ResponseWriter, r *http.Request) {
	var email string = r.FormValue("login-email")
	var password string = r.FormValue("login-password")
	existingUserID := make(chan int, 1)
	existingUserPassword := make(chan string, 1)
	existingUserError := make(chan error, 1)

	if len(email) < 7 {
		utils.AlertError(http.StatusBadRequest, w, "Email should be 8 characters or more.")
		return
	}

	if len(password) < 7 {
		utils.AlertError(http.StatusBadRequest, w, "Password should be 8 characters or more.")
		return
	}

	go func() {
		var userID int
		var userPassword string

		err := a.db.QueryRow("select id, password from users where email = $1", email).Scan(&userID, &userPassword)

		if err != nil {
			existingUserError <- err
			close(existingUserID)
			close(existingUserError)
			close(existingUserPassword)

			return
		}

		existingUserID <- userID
		existingUserPassword <- userPassword

		close(existingUserID)
		close(existingUserError)
		close(existingUserPassword)
	}()

	select {
	case userID := <-existingUserID:
		pw := <-existingUserPassword
		err := bcrypt.CompareHashAndPassword([]byte(pw), []byte(password))

		if err != nil {
			utils.AlertError(http.StatusNotFound, w, "User not existing.")
			return
		}

		utils.SetCookie(w, userID)

	case err := <-existingUserError:
		fmt.Println(err)
		utils.AlertError(http.StatusNotFound, w, "User not existing.")
		return
	}
}

func (a *auth) register(w http.ResponseWriter, r *http.Request) {
	var email string = r.FormValue("email")
	var name string = strings.TrimSpace(r.FormValue("name"))
	var password string = r.FormValue("password")

	startTime := time.Now()
	hashedPassword := make(chan string, 1)
	hashError := make(chan error, 1)
	emailExistingID := make(chan int, 1)

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
		user := new(models.User)

		a.db.QueryRow("select id from users where email = $1", email).Scan(&user.ID)

		emailExistingID <- user.ID

		close(emailExistingID)
	}()

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
		existingID := <-emailExistingID

		if existingID != 0 {
			utils.AlertError(http.StatusConflict, w, "Email already existing.")
			return
		}

		var id int

		var query string = `insert into users (createdAt, email, image, name, password, updatedAt)
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

		utils.SetCookie(w, id)

		fmt.Printf("\n Time Since: %s \n", time.Since(startTime))

		utils.AlertSuccess(201, w, "Registered successfully.")
	case he := <-hashError:
		if he != nil {
			fmt.Printf("\nError: %v\n", he)
			utils.AlertError(http.StatusInternalServerError, w, "Internal Server Error")
			return
		}
	}
}

func (a *auth) logout(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:     "session_token",
		HttpOnly: true,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
	}

	http.SetCookie(w, cookie)
	w.WriteHeader(200)
}
