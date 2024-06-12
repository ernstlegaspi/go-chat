package ws

import "fmt"

type Room struct {
	users map[*User]bool
	user  chan *User
}

func NewRoom() *Room {
	return &Room{
		users: make(map[*User]bool),
		user:  make(chan *User),
	}
}

func (r *Room) CreateRoom() {
	for {
		select {
		case user := <-r.user:
			r.users[user] = true
			fmt.Printf("Connected User: %s\n", user.name)
		}
	}
}
