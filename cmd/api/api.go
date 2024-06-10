package api

import (
	"gochat/internal/db"
	"gochat/internal/handlers"
	"gochat/internal/utils"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
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

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")

		if err != nil || cookie == nil {
			utils.ExecHTML(w, "index.html")
			return
		}

		token, tokenErr := utils.ParseJWT(cookie.Value)

		if tokenErr != nil {
			utils.ExecHTML(w, "index.html")
			return
		}

		tokenID := token.Claims.(jwt.MapClaims)["id"]

		if tokenID == 0 {
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
