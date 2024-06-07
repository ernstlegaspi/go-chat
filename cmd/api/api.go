package api

import (
	"html/template"
	"net/http"

	"gochat/internal/handlers"
)

type server struct {
	addr string
}

func InitServer(addr string) *server {
	return &server{
		addr: addr,
	}
}

func (s *server) InitAPI() error {
	router := http.NewServeMux()
	fs := http.FileServer(http.Dir("../internal/static"))

	router.Handle("/static/", http.StripPrefix("/static/", fs))

	auth := handlers.InitAuth(router)
	auth.InitAuthAPI()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		templ := template.Must(template.ParseFiles("../internal/views/index.html"))
		templ.Execute(w, nil)
	})

	sv := &http.Server{
		Addr:    s.addr,
		Handler: router,
	}

	return sv.ListenAndServe()
}
