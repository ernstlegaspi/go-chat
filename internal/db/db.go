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
	conn := "user=postgres dbname=postgres password=inventory-management sslmode=disable"

	dbConn, err := sql.Open("postgres", conn)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	defer dbConn.Close()

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
		name varchar(50),
		email varchar(100),
		password text,
		updatedAt timestamp
	)`

	_, err := d.db.Exec(query)

	return err
}
