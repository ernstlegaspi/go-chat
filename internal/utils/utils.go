package utils

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func AlertError(status int, w http.ResponseWriter, message string) {
	w.Header().Set("HX-Target", "#error-message")
	w.WriteHeader(status)
	templ := template.Must(template.ParseFiles("../internal/views/alerts/error.html"))
	templ.Execute(w, map[string]string{"message": message})
}

func AlertSuccess(w http.ResponseWriter, message string) {
	w.Header().Set("HX-Target", "#success-message")
	templ := template.Must(template.ParseFiles("../internal/views/alerts/success.html"))
	templ.Execute(w, map[string]string{"message": message})
}

func LettersAndSpaces(str string) bool {
	regex, _ := regexp.Compile("^[a-zA-Z ]+$")
	return regex.MatchString(str)
}

func CreateJWT(id int) string {
	expiration := time.Second * time.Duration(3600*24*7)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":      id,
		"expires": expiration,
	})

	str, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	return str
}
