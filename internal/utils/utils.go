package utils

import (
	"html/template"
	"net/http"
	"regexp"
)

func AlertError(w http.ResponseWriter, message string) {
	w.Header().Set("HX-Target", "#error-message")
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
