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
	w.WriteHeader(status)
	templ := template.Must(template.ParseFiles("../internal/views/alerts/error.html"))
	templ.Execute(w, map[string]string{"message": message})
}

func AlertSuccess(status int, w http.ResponseWriter, message string) {
	w.Header().Set("HX-Target", "#success-message")
	w.WriteHeader(status)
	templ := template.Must(template.ParseFiles("../internal/views/alerts/success.html"))
	templ.Execute(w, map[string]string{"message": message})
}

func LettersAndSpaces(str string) bool {
	regex, _ := regexp.Compile("^[a-zA-Z ]+$")
	return regex.MatchString(str)
}

func HasJWT(r *http.Request) (jwt.MapClaims, error) {
	cookie, err := r.Cookie("session_token")

	if err != nil || cookie == nil {
		return nil, err
	}

	token, parseErr := ParseJWT(cookie.Value)

	if parseErr != nil {
		return nil, parseErr
	}

	mapClaims := token.Claims.(jwt.MapClaims)

	return mapClaims, nil
}

func ParseJWT(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("token error")
		}

		return []byte(os.Getenv("SECRET_KEY")), nil
	})
}

func SetCookie(w http.ResponseWriter, id int, name string) {
	expiration := time.Second * time.Duration(3600*24*7)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":      id,
		"name":    name,
		"expires": expiration,
	})

	fmt.Println(os.Getenv("SECRET_KEY"))

	tokenStr, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	cookie := &http.Cookie{
		Name:     "session_token",
		HttpOnly: true,
		Value:    tokenStr,
		Path:     "/",
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	}

	http.SetCookie(w, cookie)
}

func ExecHTML(w http.ResponseWriter, path string) {
	templ := template.Must(template.ParseFiles("../internal/views/" + path))
	templ.Execute(w, nil)
}
