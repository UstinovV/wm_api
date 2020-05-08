package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	)

type DB struct {
	connection string
	db *sql.DB
}

func NewDB(connection string) *DB{
	return &DB{
		connection: connection,
	}
}

func (DB *DB) Open() error {
	fmt.Println(DB.connection)
	db, err := sql.Open("postgres", DB.connection)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	DB.db = db

	return nil
}

func (DB *DB) Close() {
	DB.db.Close()
}



