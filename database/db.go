package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	)

type Database struct {
	connection string
	db *sql.DB
}

func New(connection string) *Database{
	return &Database{
		connection: connection,
	}
}

func (database *Database) Open() error {
	fmt.Println(database.connection)
	db, err := sql.Open("postgres", database.connection)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	database.db = db

	return nil
}

func (database *Database) Close() {
	database.db.Close()
}



