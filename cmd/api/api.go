package api

import (
	"gochat/internal/db"
	"gochat/internal/handlers"
	"gochat/internal/utils"
	"net/http"
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

	dbs, database := db.CreateDB()

	dbs.CreateTables()

	router.Handle("/static/", http.StripPrefix("/static/", fs))

	auth := handlers.InitAuth(database, router)
	auth.InitAuthAPI()

	user := handlers.InitUser(database, router)
	user.InitUserAPI()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		claims, err := utils.HasJWT(r)
		tokenID := claims["id"]

		if err != nil || tokenID == 0 {
			utils.ExecHTML(w, "index.html")
			return
		}

		var userID int

		queryErr := database.QueryRow("select id from users where id = $1", tokenID).Scan(&userID)

		if queryErr != nil || userID == 0 {
			utils.ExecHTML(w, "index.html")
			return
		}

		utils.ExecHTML(w, "home.html")
	})

	router.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		utils.ExecHTML(w, "home.html")
	})

	sv := &http.Server{
		Addr:    s.addr,
		Handler: router,
	}

	return sv.ListenAndServe()
}
