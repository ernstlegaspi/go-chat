package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type db struct {
	db *sql.DB
}

func CreateDB() (*db, *sql.DB) {
	conn := "user=postgres dbname=postgres password=gochat sslmode=disable"

	dbConn, err := sql.Open("postgres", conn)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	return &db{
		db: dbConn,
	}, dbConn
}

func (d *db) CreateTables() {
	if err := d.createUserTable(); err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func (d *db) createUserTable() error {
	query := `create table if not exists users (
		id serial primary key,
		createdAt timestamp,
		email varchar(100),
		image text,
		name varchar(50),
		password text,
		updatedAt timestamp
	)`

	_, err := d.db.Exec(query)

	return err
}
