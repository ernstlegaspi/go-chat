package models

import "time"

type User struct {
	ID        int
	CreatedAt time.Time
	Email     string
	Image     string
	Name      string
	Password  string
	UpdatedAt time.Time
}
