package main

import (
	_ "github.com/lib/pq"
	"database/sql"
	"fmt"
)

const (
	host     = "172.20.0.3"
	port     = 5432
	user     = "postgres"
	password = ""
	dbname   = "wm"
)

func main() {
	fmt.Println("API")
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()
}