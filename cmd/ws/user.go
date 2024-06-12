package ws

import (
	"gochat/internal/utils"
	"net/http"

	"github.com/gorilla/websocket"
)

type User struct {
	conn *websocket.Conn
	chat chan []byte
	name string
	room *Room
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func ConnectRoom(room *Room, w http.ResponseWriter, r *http.Request) {
	conn, connErr := upgrader.Upgrade(w, r, nil)

	if connErr != nil {
		utils.AlertError(500, w, "Internal server error. Try again later.")
		return
	}

	claims, claimsErr := utils.HasJWT(r)
	userName, ok := claims["name"].(string)

	if !ok {
		utils.AlertError(500, w, "Internal server error. Login in again.")
		return
	}

	if claimsErr != nil {
		utils.AlertError(400, w, "Unable to process request. Try again later.")
		return
	}

	user := &User{
		conn: conn,
		chat: make(chan []byte),
		name: userName,
		room: room,
	}

	room.user <- user
}
