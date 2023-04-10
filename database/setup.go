package database

import (
	"bookstore/helper"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "mukul98s"
	password = "Mukul@14"
	dbname   = "bookstore"
)

var DB *sql.DB = setupDB()

func setupDB() *sql.DB {

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlconn)
	helper.CheckError(err, "DB Connection Failed....")

	// assigning new error to the pre-existing variable
	err = db.Ping()
	helper.CheckError(err, "DB Ping Failed....")

	return db
}
