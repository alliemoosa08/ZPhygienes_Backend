package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "P@ssword08"
	dbname   = "zpHygieneDB"
)

func main() {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)

	defer db.Close()

	//check if Database is connected.
	err = db.Ping()
	CheckError(err)

	fmt.Println("Database connected")
}
func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
