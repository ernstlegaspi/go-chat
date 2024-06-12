package handlers

import (
	"database/sql"
	"fmt"
	"gochat/internal/models"
	"gochat/internal/utils"
	"html/template"
	"net/http"
	"strconv"
)

type users struct {
	db *sql.DB
	r  *http.ServeMux
}

func InitUser(db *sql.DB, r *http.ServeMux) *users {
	return &users{
		db: db,
		r:  r,
	}
}

func (u *users) InitUserAPI() {
	var r *http.ServeMux = u.r

	r.HandleFunc("GET /users", u.getUsers)
}

func (u *users) getUsers(w http.ResponseWriter, r *http.Request) {
	claims, claimsErr := utils.HasJWT(r)

	if claimsErr != nil {
		fmt.Println(claimsErr)
		utils.AlertError(400, w, "Can not process request. Try again later.")
		return
	}

	userID := claims["id"]

	rows, err := u.db.Query("select id, image, name from users where id != $1", userID)

	if err != nil {
		fmt.Println(err)
		utils.AlertError(500, w, "Internal Server Error")
		return
	}

	defer rows.Close()

	for rows.Next() {
		user := new(models.User)

		err := rows.Scan(
			&user.ID,
			&user.Image,
			&user.Name,
		)

		if err != nil {
			fmt.Println(err)
			return
		}

		templ := template.Must(template.ParseFiles("../internal/views/user.html"))
		templ.ExecuteTemplate(w, "user-card", map[string]string{
			"ID":    strconv.Itoa(user.ID),
			"Image": user.Image,
			"Name":  user.Name,
		})
	}
}
