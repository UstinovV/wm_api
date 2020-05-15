package database

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

type DB struct {
	connection string
	Db         *sql.DB
}

func NewDB(connection string) *DB {
	return &DB{
		connection: connection,
	}
}

func (DB *DB) Open() error {
	db, err := sql.Open("postgres", DB.connection)
	if err != nil {
		log.Fatal("DB open: ", err)
		return err
	}

	if err := db.Ping(); err != nil {
		log.Fatal("DB Ping: ", err)
		return err
	}

	DB.Db = db

	return nil
}

func (DB *DB) Close() {
	DB.Db.Close()
}
